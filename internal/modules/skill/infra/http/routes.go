package http

import (
	"github.com/boilerplate/internal/modules/skill/application"
	"github.com/boilerplate/internal/modules/skill/infra/persistance"
	"github.com/boilerplate/internal/shared/middleware"
	"github.com/boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, db *pgxpool.Pool, log logger.Logger) {
	repo := persistance.NewSkillRepoImpl(db)
	service := application.NewSkillService(repo)
	handler := NewSkillHandler(service)

	group := e.Group("/skills")
	group.GET("", handler.List)
	group.GET("/grouped", handler.GroupedByCategory)
	group.Use(middleware.Auth(), middleware.IsAdmin())
	group.POST("", handler.Create)
	group.GET("/:id", handler.GetByID)
	group.PATCH("/:id", handler.Update)
	group.DELETE("/:id", handler.Delete)
}
