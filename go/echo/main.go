package main

import (
	"context"
	"fmt"
	"net/http"
	"shortr/cache"
	"shortr/config"
	"shortr/repo"
	"shortr/shortid"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var idCache = cache.New(1024)
var db *pgxpool.Pool

func getURL(ctx echo.Context) error {
	sid := ctx.Param("id")
	if value, exists := idCache.Read(sid); exists {
		return ctx.String(http.StatusOK, fmt.Sprintf("CACHED: %s", value))
	}
	sint, _ := strconv.Atoi(sid)
	id, _ := shortid.Encode(sint)
	idCache.Write(sid, id)

	return ctx.String(http.StatusOK, fmt.Sprintf("%s", id))
}

func shortenURL(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Shorten URL endpoint!")
}

func deleteURL(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Delete URL endpoint!")
}

func modifyURL(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Modify URL endpoint!")
}

func getURLStats(ctx echo.Context) error {
	return ctx.String(http.StatusOK, "Get URL stats endpoint!")
}

func main() {
	db, err := repo.Connect(context.Background(), fmt.Sprintf("postgresql://%s:%s@%s:%d",
		config.GetEnvAsString("DATABASE_USER", "postgres"),
		config.GetEnvAsString("DATABASE_PASSWORD", "postgres"),
		config.GetEnvAsString("DATABASE_HOST", "postgres"),
		config.GetEnvAsInt("DATABASE_PORT", 5432),
	), 5)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := echo.New()
	app.Use(middleware.Logger())

	// Routes
	app.GET("/:id", getURL)
	app.POST("/:id", shortenURL)
	app.DELETE("/:id", deleteURL)
	app.PUT("/:id", modifyURL)
	app.GET("/:id/stats", getURLStats)

	app.Logger.Fatal(app.Start(fmt.Sprintf(":%d", config.GetEnvAsInt("APP_PORT", 80))))
}
