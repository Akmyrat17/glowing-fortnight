package server

import (
	"fmt"

	"github.com/boilerplate/internal/config"
	shared_middleware "github.com/boilerplate/internal/shared/middleware"
	"github.com/boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New(cfg *config.Config, db *pgxpool.Pool, log logger.Logger) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(shared_middleware.ErrorMiddlewareWithLogger(log))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))
	// Serve uploaded files
	e.Static("/uploads", "uploads")

	registerRoutes(e, db, cfg, log)
	return e
}

func Start(e *echo.Echo, port int) error {
	return e.Start(fmt.Sprintf(":%d", port))
}
