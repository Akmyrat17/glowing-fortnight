package http

import (
	"net/http"
	"sort"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/skill/application"
	"github.com/boilerplate/internal/modules/skill/infra/http/dto"
	"github.com/boilerplate/internal/shared/response"
	"github.com/boilerplate/pkg/query"
	reqctx "github.com/boilerplate/pkg/req_ctx"
	"github.com/labstack/echo/v4"
)

type SkillHandler struct {
	service *application.SkillService
}

func NewSkillHandler(service *application.SkillService) *SkillHandler {
	return &SkillHandler{service: service}
}

func (h *SkillHandler) Create(c echo.Context) error {
	var req dto.CreateSkillReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}

	skill, err := req.ToDomain()
	if err != nil {
		return err
	}

	if err := h.service.Create(c.Request().Context(), skill); err != nil {
		return err
	}

	return response.Created(c, dto.SkillResFromDomain(skill))
}

func (h *SkillHandler) GetByID(c echo.Context) error {
	id, err := domain.ParseSkillID(c.Param("id"))
	if err != nil {
		return err
	}

	skill, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return response.OK(c, dto.SkillResFromDomain(skill))
}

func (h *SkillHandler) Update(c echo.Context) error {
	id, err := domain.ParseSkillID(c.Param("id"))
	if err != nil {
		return err
	}

	var req dto.UpdateSkillReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}

	existing, err := h.service.FindByID(c.Request().Context(), id)
	if err != nil {
		return err
	}

	skill, err := req.ToDomain(existing)
	if err != nil {
		return err
	}
	skill.ID = id

	if err := h.service.Update(c.Request().Context(), skill); err != nil {
		return err
	}

	return response.OK(c, dto.SkillResFromDomain(skill))
}

func (h *SkillHandler) Delete(c echo.Context) error {
	id, err := domain.ParseSkillID(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *SkillHandler) List(c echo.Context) error {
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

	return response.Paginated(c, dto.SkillResFromDomainList(list), pagination.Page, pagination.PerPage, total)
}

func (h *SkillHandler) GroupedByCategory(c echo.Context) error {
	groups, err := h.service.FindAllGroupedByCategory(c.Request().Context())
	if err != nil {
		return err
	}

	categories := make([]domain.SkillCategory, 0, len(groups))
	for category := range groups {
		categories = append(categories, category)
	}
	sort.Slice(categories, func(i, j int) bool {
		return categories[i] < categories[j]
	})

	res := make([]*dto.SkillCategoryGroupRes, 0, len(categories))
	for _, category := range categories {
		res = append(res, &dto.SkillCategoryGroupRes{
			Category: string(category),
			Skills:   dto.SkillResFromDomainList(groups[category]),
		})
	}

	return response.OK(c, res)
}
