package main

import (
	"fmt"
	"net/http"
	nurl "net/url"
	"shortr/cache"
	"shortr/config"
	"shortr/render"
	"shortr/repo"
	"shortr/shortid"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var urlCache = cache.New(1024)
var urlRepo *repo.Repo

func getURL(ctx echo.Context) error {
	name := ctx.Param("name")

	if url, exists := urlCache.Read(name); exists {
		go logIfErr(ctx.Logger(), wrap(urlRepo.UpdateMetricsByName(name))...)
		return ctx.Redirect(http.StatusTemporaryRedirect, url.(string)) // HTTP CODE 307 IN ORDER NOT TO GET URLs CACHED
	}

	url, err := urlRepo.GetByName(name)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrNotFound
	}

	go urlCache.Write(*url.Name, url.URL)
	go logIfErr(ctx.Logger(), wrap(urlRepo.UpdateMetricsByName(name))...)

	return ctx.Redirect(http.StatusTemporaryRedirect, url.URL) // HTTP CODE 307 IN ORDER NOT TO GET URLs CACHED
}

func shortenURL(ctx echo.Context) error {
	name := ctx.Param("name")
	qurl := ctx.QueryParam("url")

	_, err := nurl.ParseRequestURI(qurl)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrBadRequest
	}

	url, err := urlRepo.Create(qurl)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	if name == "" {
		name, err = shortid.Encode(url.ID)
	}
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	url, err = urlRepo.UpdateNameByID(url.ID, name)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusOK, url)
}

func deleteURL(ctx echo.Context) error {
	name := ctx.Param("name")

	url, err := urlRepo.DeleteByName(name)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	urlCache.Remove(*url.Name)

	return ctx.JSON(http.StatusOK, url)
}

func modifyURL(ctx echo.Context) error {
	name := ctx.Param("name")
	qurl := ctx.QueryParam("url")

	_, err := nurl.ParseRequestURI(qurl)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrBadRequest
	}

	url, err := urlRepo.UpdateURLByName(name, qurl)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	if _, exists := urlCache.Read(*url.Name); exists {
		urlCache.Write(*url.Name, url.URL)
	}

	return ctx.JSON(http.StatusOK, url)
}

func getURLStats(ctx echo.Context) error {
	name := ctx.Param("name")
	contentType := ctx.Request().Header.Get(echo.HeaderContentType)

	url, err := urlRepo.GetByName(name)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrNotFound
	}

	switch contentType {
	case echo.MIMEApplicationJSON, echo.MIMEApplicationJSONCharsetUTF8:
		return ctx.JSON(http.StatusOK, url)
	default:
		return ctx.Render(http.StatusOK, "stats.gts.html", url)
	}
}

func main() {
	var err error
	urlRepo, err = repo.Connect(fmt.Sprintf("postgresql://%s:%s@%s:%d",
		config.GetEnvAsString("DATABASE_USER", "postgres"),
		config.GetEnvAsString("DATABASE_PASSWORD", "postgres"),
		config.GetEnvAsString("DATABASE_HOST", "postgres"),
		config.GetEnvAsInt("DATABASE_PORT", 5432),
	), 5)
	if err != nil {
		panic(err)
	}
	defer urlRepo.Disconnect()

	scheme := "http"
	if config.GetEnvAsBool("APP_SSL_ENABLED", false) {
		scheme = "https"
	}

	app := echo.New()
	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(middleware.Logger()) // TODO : Find another logger, this is too slow
	app.Use(middleware.Recover())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{fmt.Sprintf("%s://%s",
			scheme,
			config.GetEnvAsSlice("VIRTUAL_HOST", []string{"localhost"})[0]),
		}, // Add more origins if required
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
	}))
	app.HTTPErrorHandler = customHTTPErrorHandler
	app.Renderer = render.New("/static/templates/*.gts.html")
	app.IPExtractor = echo.ExtractIPFromRealIPHeader()

	// Routes
	app.Static("/", "/static")
	app.POST("/", shortenURL)
	url := app.Group("/:name")
	/*--*/ url.GET("", getURL)
	/*--*/ url.POST("", shortenURL)
	/*--*/ url.DELETE("", deleteURL)
	/*--*/ url.PUT("", modifyURL)
	/*--*/ url.GET("/stats", getURLStats)

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", config.GetEnvAsInt("APP_PORT", 80))))
}

func customHTTPErrorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	if httpError, ok := err.(*echo.HTTPError); ok {
		code = httpError.Code
	}

	contentType := ctx.Request().Header.Get(echo.HeaderContentType)
	switch contentType {
	case echo.MIMEApplicationJSON, echo.MIMEApplicationJSONCharsetUTF8:
		ctx.Echo().DefaultHTTPErrorHandler(err, ctx)
		return
	}
	ctx.File(fmt.Sprintf("/static/templates/%d.html", code))
	ctx.Echo().DefaultHTTPErrorHandler(err, ctx)
}

func wrap(vs ...interface{}) []interface{} {
	return vs
}

func logIfErr(logger echo.Logger, i ...interface{}) {
	var err error = nil
	for _, o := range i {
		switch t := o.(type) {
		case error:
			err = t
		}
	}
	if err != nil {
		logger.Error(err)
	}
}
