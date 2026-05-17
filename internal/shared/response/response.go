package response

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    *PageMeta   `json:"meta"`
}

type PageMeta struct {
	Page      int   `json:"page"`
	PerPage   int   `json:"per_page"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_pages"`
}

type ErrorResponse struct {
	Type      string `json:"type"`
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	Timestamp string `json:"timestamp"`
}

func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func Paginated(c echo.Context, data interface{}, page, perPage int, total int64) error {
	totalPages := (total + int64(perPage) - 1) / int64(perPage)
	return c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    data,
		Meta: &PageMeta{
			Page:      page,
			PerPage:   perPage,
			Total:     total,
			TotalPage: totalPages,
		},
	})
}

func Error(c echo.Context, errorType string, status int, message string) error {
	return c.JSON(status, ErrorResponse{
		Type:      errorType,
		Status:    status,
		Message:   message,
		Path:      c.Request().URL.Path,
		Method:    c.Request().Method,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
