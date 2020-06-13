package main

import (
	"fmt"
	"net/http"
	nurl "net/url"
	"shortr/cache"
	"shortr/config"
	"shortr/repo"
	"shortr/shortid"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var urlCache = cache.New(1024)
var urlRepo *repo.Repo

func getURL(ctx echo.Context) error {
	name := ctx.Param("name")

	if value, exists := urlCache.Read(name); exists {
		go logIfErr(ctx, urlRepo.UpdateMetricsByName(name))
		return ctx.Redirect(http.StatusMovedPermanently, value.(string))
	}

	url, err := urlRepo.GetByName(name)
	if err != nil {
		// TODO : 404 NOT FOUND PRETTY HTML
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	go urlCache.Write(name, url)
	go logIfErr(ctx, urlRepo.UpdateMetricsByName(name))

	return ctx.Redirect(http.StatusMovedPermanently, url)
}

func shortenURL(ctx echo.Context) error {
	name := ctx.Param("name")
	url := ctx.QueryParam("url")

	_, err := nurl.ParseRequestURI(url)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	id, err := urlRepo.Create(url)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	if name == "" {
		name, err = shortid.Encode(id)
	}
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err = urlRepo.UpdateNameByID(id, name)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	return ctx.String(http.StatusOK, fmt.Sprintf("Created %s with %s\n", name, url))
}

func deleteURL(ctx echo.Context) error {
	name := ctx.Param("name")

	err := urlRepo.DeleteByName(name)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	urlCache.Remove(name)

	return ctx.String(http.StatusOK, fmt.Sprintf("Deleted %s\n", name))
}

func modifyURL(ctx echo.Context) error {
	name := ctx.Param("name")
	url := ctx.QueryParam("url")

	_, err := nurl.ParseRequestURI(url)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	err = urlRepo.UpdateURLByName(name, url)
	if err != nil {
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	if _, exists := urlCache.Read(name); exists {
		urlCache.Write(name, url)
	}

	return ctx.String(http.StatusOK, fmt.Sprintf("Updated %s to %s\n", name, url))
}

func getURLStats(ctx echo.Context) error {
	name := ctx.Param("name")
	contentType := ctx.Request().Header.Get("Content-Type")

	hits, lastHitAt, createdAt, modifiedAt, err := urlRepo.GetMetricsByName(name)
	if err != nil {
		// TODO : 404 NOT FOUND PRETTY HTML
		ctx.Logger().Error(err)
		return ctx.NoContent(http.StatusBadRequest)
	}

	switch contentType {
	case "application/json":
		var response struct {
			Hits       int        `json:"hits"`
			LastHitAt  *time.Time `json:"last_hit_at"`
			CreatedAt  time.Time  `json:"created_at"`
			ModifiedAt time.Time  `json:"modified_at"`
		}
		response.Hits = hits
		response.LastHitAt = lastHitAt
		response.CreatedAt = createdAt
		response.ModifiedAt = modifiedAt
		return ctx.JSONPretty(http.StatusOK, response, "  ")
	default:
		// TODO : RENDER PRETTY HTML IF Content-Type: text/html or default
		return ctx.String(http.StatusOK, "Get URL stats endpoint!")
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

	app := echo.New()
	app.Use(middleware.Logger()) // TODO : Find another logger, this is too slow

	// Routes
	// app.Static("/", "/assets/index.html") // TODO : RENDER PRETTY HTML TO INTERACT WITH THE API
	app.POST("/", shortenURL)
	url := app.Group("/:name")
	/*--*/ url.GET("", getURL)
	/*--*/ url.POST("", shortenURL)
	/*--*/ url.DELETE("", deleteURL)
	/*--*/ url.PUT("", modifyURL)
	/*--*/ url.GET("/stats", getURLStats)

	// TODO : RENDER PRETTY HTML WHEN ERROR

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", config.GetEnvAsInt("APP_PORT", 80))))
}

func logIfErr(ctx echo.Context, err error) {
	if err != nil {
		ctx.Logger().Error(err)
	}
}
