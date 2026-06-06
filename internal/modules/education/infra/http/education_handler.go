package http

import (
	"net/http"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/education/application"
	"github.com/boilerplate/internal/modules/education/infra/http/dto"
	"github.com/boilerplate/internal/shared/response"
	"github.com/boilerplate/pkg/query"
	reqctx "github.com/boilerplate/pkg/req_ctx"
	"github.com/labstack/echo/v4"
)

type EducationHandler struct {
	service *application.EducationService
}

func NewEducationHandler(service *application.EducationService) *EducationHandler {
	return &EducationHandler{service: service}
}

func (h *EducationHandler) Create(c echo.Context) error {
	var req dto.CreateEducationReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}

	e := req.ToDomain()
	if err := h.service.Create(c.Request().Context(), e); err != nil {
		return err
	}

	return response.Created(c, dto.EducationResFromDomain(e))
}

func (h *EducationHandler) GetByID(c echo.Context) error {
	id, err := domain.ParseEducationID(c.Param("id"))
	if err != nil {
		return err
	}

	e, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return response.OK(c, dto.EducationResFromDomain(e))
}

func (h *EducationHandler) Update(c echo.Context) error {
	id, err := domain.ParseEducationID(c.Param("id"))
	if err != nil {
		return err
	}

	var req dto.UpdateEducationReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}

	existing, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	e := req.ToDomain(existing)
	e.ID = id

	if err := h.service.Update(c.Request().Context(), e); err != nil {
		return err
	}

	return response.OK(c, dto.EducationResFromDomain(e))
}

func (h *EducationHandler) Delete(c echo.Context) error {
	id, err := domain.ParseEducationID(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *EducationHandler) List(c echo.Context) error {
	pagination := reqctx.ParsePagination(c)
	filters := query.ParseFilters(c.QueryParams())
	sorts := query.ParseSort(c.QueryParam("sort"))

	list, total, err := h.service.FindAll(
		c.Request().Context(),
		pagination.PerPage,
		pagination.Offset(),
		filters,
		sorts,
	)
	if err != nil {
		return err
	}

	return response.Paginated(c, dto.EducationResFromDomainList(list), pagination.Page, pagination.PerPage, total)
}
