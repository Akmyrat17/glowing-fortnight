package http

import (
	"github.com/boilerplate/internal/config"
	"github.com/boilerplate/internal/modules/upload/application"
	"github.com/boilerplate/pkg/logger"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Group, log logger.Logger, cfg config.Cloudinary) {
	uploadService := application.NewUploadService(log, cfg)
	handler := NewUploadHandler(uploadService)
	group := e.Group("/uploads")
	group.POST("/image", handler.UploadImage)
}
