package server

import (
	"github.com/boilerplate/internal/config"
	authhttp "github.com/boilerplate/internal/modules/auth/infra/http"
	educationhttp "github.com/boilerplate/internal/modules/education/infra/http"
	experiencehttp "github.com/boilerplate/internal/modules/experience/infra/http"
	permissionhttp "github.com/boilerplate/internal/modules/permission/infra/http"
	profilehttp "github.com/boilerplate/internal/modules/profile/infra/http"
	projecthttp "github.com/boilerplate/internal/modules/projects/infra/http"
	skillhttp "github.com/boilerplate/internal/modules/skill/infra/http"
	uploadhttp "github.com/boilerplate/internal/modules/upload/infra/http"
	userhttp "github.com/boilerplate/internal/modules/user/infra/http"
	"github.com/boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func registerRoutes(e *echo.Echo, db *pgxpool.Pool, cfg *config.Config, log logger.Logger) {
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	e.Static("/uploads", "uploads")

	api := e.Group("/api/v1")
	authhttp.RegisterRoutes(api, db, cfg, log)
	userhttp.RegisterRoutes(api, db, log)
	permissionhttp.RegisterRoutes(api, db, log)
	profilehttp.RegisterRoutes(api, db, log)
	uploadhttp.RegisterRoutes(api, log)
	experiencehttp.RegisterRoutes(api, db, log)
	educationhttp.RegisterRoutes(api, db, log)
	skillhttp.RegisterRoutes(api, db, log)
	projecthttp.RegisterRoutes(api, db, log)
}
