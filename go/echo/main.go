package main

import (
	"context"
	"fmt"
	"net/http"
	nurl "net/url"
	"os"
	"os/signal"
	"shortr/cache"
	"shortr/config"
	"shortr/logger"
	"shortr/model"
	"shortr/render"
	"shortr/repo"
	"shortr/shortid"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var urlCache = cache.New(4096)
var urlRepo *repo.Repo

func getURL(ctx echo.Context) error {
	name := ctx.Param("name")

	if url, exists := urlCache.Read(name); exists {
		go logIfErr(ctx.Logger(), wrap(urlRepo.UpdateMetricsByName(ctx.Request().Context(), name))...)
		return ctx.Redirect(http.StatusTemporaryRedirect, url.(string)) // HTTP CODE 307 IN ORDER NOT TO GET URLs CACHED
	}

	url, err := urlRepo.GetByName(ctx.Request().Context(), name)
	if err != nil {
		if err == repo.ErrNoRows {
			return echo.ErrNotFound
		}
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	go urlCache.Write(url.Name, url.URL)
	go logIfErr(ctx.Logger(), wrap(urlRepo.UpdateMetricsByName(ctx.Request().Context(), name))...)

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

	var url model.URL
	err = urlRepo.Transaction(ctx.Request().Context(), func(urlTxRepo *repo.Repo) error {
		url, err = urlTxRepo.Create(ctx.Request().Context(), qurl)
		if err != nil {
			return err
		}

		if name == "" {
			name, err = shortid.Encode(url.ID)
		}
		if err != nil {
			return err
		}

		url, err = urlTxRepo.UpdateNameByID(ctx.Request().Context(), url.ID, name)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		if err == repo.ErrIntegrityViolation {
			return echo.ErrBadRequest
		}
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	return ctx.JSON(http.StatusOK, url)
}

func deleteURL(ctx echo.Context) error {
	name := ctx.Param("name")

	url, err := urlRepo.DeleteByName(ctx.Request().Context(), name)
	if err != nil {
		if err == repo.ErrNoRows {
			return echo.ErrBadRequest
		}
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	urlCache.Remove(url.Name)

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

	url, err := urlRepo.UpdateURLByName(ctx.Request().Context(), name, qurl)
	if err != nil {
		if err == repo.ErrNoRows || err == repo.ErrIntegrityViolation {
			return echo.ErrBadRequest
		}
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
	}

	if _, exists := urlCache.Read(url.Name); exists {
		urlCache.Write(url.Name, url.URL)
	}

	return ctx.JSON(http.StatusOK, url)
}

func getURLStats(ctx echo.Context) error {
	name := ctx.Param("name")
	contentType := ctx.Request().Header.Get(echo.HeaderContentType)

	url, err := urlRepo.GetByName(ctx.Request().Context(), name)
	if err != nil {
		if err == repo.ErrNoRows {
			return echo.ErrNotFound
		}
		ctx.Logger().Error(err)
		return echo.ErrInternalServerError
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
	appLogger := logger.New("shortr")

	urlRepo, err = repo.Connect(fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		config.GetEnvAsString("DATABASE_USER", "postgres"),
		config.GetEnvAsString("DATABASE_PASSWORD", "postgres"),
		config.GetEnvAsString("DATABASE_HOST", "postgres"),
		config.GetEnvAsInt("DATABASE_PORT", 5432),
		config.GetEnvAsString("DATABASE_NAME", "postgres"),
		config.GetEnvAsString("DATABASE_SSLMODE", "disable"),
	), 5, logger.Database(appLogger))
	if err != nil {
		panic(err)
	}
	defer urlRepo.Disconnect()

	scheme := "http"
	if config.GetEnvAsBool("APP_SSL_ENABLED", false) {
		scheme = "https"
	}

	app := echo.New()
	app.Logger = logger.Standard(appLogger)
	app.HTTPErrorHandler = customHTTPErrorHandler
	app.Renderer = render.New("/static/templates/*.gts.html")
	app.IPExtractor = echo.ExtractIPFromRealIPHeader()

	app.Pre(middleware.RemoveTrailingSlash())
	app.Use(logger.Middleware(appLogger))
	app.Use(middleware.Recover())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{fmt.Sprintf("%s://%s",
			scheme,
			config.GetEnvAsSlice("VIRTUAL_HOST", []string{"localhost"})[0]),
		}, // Add more origins if required
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
	}))

	// Routes
	app.Static("/", "/static")
	app.GET("/health", healthCheck)
	app.POST("/", shortenURL)
	url := app.Group("/:name")
	/*--*/ url.GET("", getURL)
	/*--*/ url.POST("", shortenURL)
	/*--*/ url.DELETE("", deleteURL)
	/*--*/ url.PUT("", modifyURL)
	/*--*/ url.GET("/stats", getURLStats)

	go app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", config.GetEnvAsInt("APP_PORT", 80))))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		app.Logger.Fatal(err)
	}
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

func healthCheck(ctx echo.Context) error {
	err := urlRepo.Health()
	if err != nil {
		ctx.Logger().Error(err)
		return echo.ErrServiceUnavailable
	}
	return ctx.String(http.StatusOK, "OK\n")
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
