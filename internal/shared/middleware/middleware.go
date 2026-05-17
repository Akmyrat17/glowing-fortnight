package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/yourorg/boilerplate/internal/domain"
	"github.com/yourorg/boilerplate/internal/shared/app_errors"
	"github.com/yourorg/boilerplate/internal/shared/response"
)

func ErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		var appErr *app_errors.AppError
		if errors.As(err, &appErr) {
			return response.Error(c, appErr.Code, appErr.Status, appErr.Message)
		}

		// Handle pgx errors
		if errors.Is(err, pgx.ErrNoRows) {
			return response.Error(c, app_errors.ErrCodeNotFound, http.StatusNotFound, "resource not found")
		}

		// Default to internal error
		return response.Error(c, app_errors.ErrCodeInternal, http.StatusInternalServerError, "an internal error occurred")
	}
}

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Request logging could go here
		start := time.Now()

		err := next(c)

		latency := time.Since(start).Milliseconds()
		_ = latency // Used in logging

		// Response logging could go here
		return err
	}
}

type TokenAuthenticator interface {
	ValidateToken(ctx context.Context, token string) (*domain.User, error)
}

var authProvider TokenAuthenticator

func SetAuthProvider(provider TokenAuthenticator) {
	authProvider = provider
}

func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return app_errors.Unauthorized("authorization header missing")
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				return app_errors.Unauthorized("invalid authorization header")
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if token == "" {
				return app_errors.Unauthorized("missing bearer token")
			}

			if authProvider == nil {
				return app_errors.InternalError("auth provider is not configured")
			}

			user, err := authProvider.ValidateToken(c.Request().Context(), token)
			if err != nil {
				return err
			}

			c.Set("user", user)
			return next(c)
		}
	}
}

func CurrentUser(c echo.Context) (*domain.User, bool) {
	user, ok := c.Get("user").(*domain.User)
	return user, ok
}

func RequirePermission(permission string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check user has permission
			// For now, just pass through
			return next(c)
		}
	}
}
