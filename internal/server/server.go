package server

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yourorg/boilerplate/internal/config"
	authhttp "github.com/yourorg/boilerplate/internal/modules/auth/infra/http"
	permissionhttp "github.com/yourorg/boilerplate/internal/modules/permission/infra/http"
	userhttp "github.com/yourorg/boilerplate/internal/modules/user/infra/http"
	shared_middleware "github.com/yourorg/boilerplate/internal/shared/middleware"
	"github.com/yourorg/boilerplate/pkg/logger"
)

func New(cfg *config.Config, db *pgxpool.Pool, log logger.Logger) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(shared_middleware.ErrorMiddleware)

	// Routes
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	authhttp.RegisterRoutes(e, db, cfg, log)
	userhttp.RegisterRoutes(e, db, log)
	permissionhttp.RegisterRoutes(e, db, log)

	return e
}

func Start(e *echo.Echo, port int) error {
	return e.Start(fmt.Sprintf(":%d", port))
}
