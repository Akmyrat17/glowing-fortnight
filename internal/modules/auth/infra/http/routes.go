package http

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/yourorg/boilerplate/internal/config"
	"github.com/yourorg/boilerplate/internal/modules/auth/application"
	userPersistence "github.com/yourorg/boilerplate/internal/modules/user/infra/persistence"
	sharedmiddleware "github.com/yourorg/boilerplate/internal/shared/middleware"
	"github.com/yourorg/boilerplate/pkg/logger"
)

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool, cfg *config.Config, log logger.Logger) {
	userRepo := userPersistence.NewUserRepoImpl(db)
	authService := application.NewAuthService(userRepo, cfg.JWT, log)
	sharedmiddleware.SetAuthProvider(authService)

	handler := NewAuthHandler(authService)
	authGroup := e.Group("auth")
	authGroup.POST("/login", handler.Login)
	authGroup.POST("/refresh", handler.Refresh)
	authGroup.POST("/logout", handler.Logout, sharedmiddleware.Auth())
}
