package http

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/yourorg/boilerplate/internal/modules/user/application"
	"github.com/yourorg/boilerplate/internal/modules/user/infra/persistence"
	"github.com/yourorg/boilerplate/internal/shared/middleware"
	"github.com/yourorg/boilerplate/pkg/logger"
)

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool, log logger.Logger) {
	repo := persistence.NewUserRepoImpl(db)
	service := application.NewUserService(repo, log)
	handler := NewUserHandler(service)

	g := e.Group("/users", middleware.Auth())
	g.GET("", handler.ListUsers)
	g.POST("", handler.CreateUser)
	g.GET("/:id", handler.GetUser)
	g.PUT("/:id", handler.UpdateUser, middleware.RequirePermission("users.update"))
	g.DELETE("/:id", handler.DeleteUser, middleware.RequirePermission("users.delete"))
}
