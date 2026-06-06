package http

import (
	"github.com/boilerplate/internal/modules/projects/application"
	"github.com/boilerplate/internal/modules/projects/infra/persistance"
	"github.com/boilerplate/internal/shared/middleware"
	"github.com/boilerplate/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, db *pgxpool.Pool, log logger.Logger) {
	projectRepo := persistance.NewProjectRepoImpl(db)
	projectService := application.NewProjectService(projectRepo)
	projectHandler := NewProjectHandler(projectService)

	group := e.Group("/projects")
	group.GET("", projectHandler.ListProjects)
	group.GET("/:id/skills", projectHandler.ListSkills)
	group.Use(middleware.Auth(), middleware.IsAdmin())
	group.POST("", projectHandler.CreateProject)
	group.GET("/:id", projectHandler.GetProject)
	group.PATCH("/:id", projectHandler.UpdateProject)
	group.DELETE("/:id", projectHandler.DeleteProject)

	// skill
	group.POST("/:id/skills", projectHandler.AddSkill)
	group.DELETE("/:id/skills", projectHandler.RemoveSkill)
}
