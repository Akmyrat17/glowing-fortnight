package http

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/yourorg/boilerplate/internal/modules/permission/application"
	"github.com/yourorg/boilerplate/internal/modules/permission/infra/http/dto"
	"github.com/yourorg/boilerplate/internal/modules/permission/infra/persistence"
	"github.com/yourorg/boilerplate/internal/shared/response"
	"github.com/yourorg/boilerplate/pkg/logger"
)

type PermissionHandler struct {
	service *application.PermissionService
}

func NewPermissionHandler(service *application.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

func (h *PermissionHandler) ListPermissions(c echo.Context) error {
	permissions, err := h.service.GetAllPermissions(c.Request().Context())
	if err != nil {
		return err
	}
	return response.OK(c, dto.PermissionListResFromDomain(permissions))
}

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool, log logger.Logger) {
	permRepo := persistence.NewPermissionRepoImpl(db)
	groupPermRepo := persistence.NewGroupPermissionRepoImpl(db)
	service := application.NewPermissionService(permRepo, groupPermRepo, log)
	handler := NewPermissionHandler(service)

	g := e.Group("/permissions")
	g.GET("", handler.ListPermissions)
}
