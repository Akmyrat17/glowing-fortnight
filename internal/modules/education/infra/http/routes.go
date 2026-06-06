package http

import (
	"github.com/boilerplate/internal/modules/education/application"
	"github.com/boilerplate/internal/modules/education/infra/persistance"
	"github.com/boilerplate/internal/shared/middleware"
	"github.com/boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, db *pgxpool.Pool, log logger.Logger) {
	repo := persistance.NewEducationRepoImpl(db)
	service := application.NewEducationService(repo)
	handler := NewEducationHandler(service)

	group := e.Group("/educations")
	group.GET("", handler.List)
	group.Use(middleware.Auth(), middleware.IsAdmin())
	group.POST("", handler.Create)
	group.GET("/:id", handler.GetByID)
	group.PATCH("/:id", handler.Update)
	group.DELETE("/:id", handler.Delete)
}
