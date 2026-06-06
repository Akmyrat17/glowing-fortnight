package http

import (
	"net/http"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/experience/application"
	"github.com/boilerplate/internal/modules/experience/infra/http/dto"
	"github.com/boilerplate/internal/shared/response"
	"github.com/boilerplate/pkg/query"
	reqctx "github.com/boilerplate/pkg/req_ctx"
	"github.com/labstack/echo/v4"
)

type ExperienceHandler struct {
	service *application.ExperienceService
}

func NewExperienceHandler(service *application.ExperienceService) *ExperienceHandler {
	return &ExperienceHandler{service: service}
}

func (h *ExperienceHandler) Create(c echo.Context) error {
	var req dto.CreateExperienceReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}

	e := req.ToDomain()
	if err := h.service.Create(c.Request().Context(), e); err != nil {
		return err
	}

	return response.Created(c, dto.ExperienceResFromDomain(e))
}

func (h *ExperienceHandler) GetByID(c echo.Context) error {
	id, err := domain.ParseExperienceID(c.Param("id"))
	if err != nil {
		return err
	}

	e, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return response.OK(c, dto.ExperienceResFromDomain(e))
}

func (h *ExperienceHandler) Update(c echo.Context) error {
	id, err := domain.ParseExperienceID(c.Param("id"))
	if err != nil {
		return err
	}

	var req dto.UpdateExperienceReq
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

	return response.OK(c, dto.ExperienceResFromDomain(e))
}

func (h *ExperienceHandler) Delete(c echo.Context) error {
	id, err := domain.ParseExperienceID(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ExperienceHandler) List(c echo.Context) error {
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

	return response.Paginated(c, dto.ExperienceResFromDomainList(list), pagination.Page, pagination.PerPage, total)
}
