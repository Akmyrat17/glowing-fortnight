package http

import (
	"net/http"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/projects/application"
	"github.com/boilerplate/internal/modules/projects/infra/http/dto"
	"github.com/boilerplate/internal/shared/response"
	"github.com/boilerplate/pkg/query"
	reqctx "github.com/boilerplate/pkg/req_ctx"
	"github.com/labstack/echo/v4"
)

type ProjectHandler struct {
	projectService *application.ProjectService
}

func NewProjectHandler(projectService *application.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) CreateProject(c echo.Context) error {
	var req dto.CreateProjectReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}

	project := req.ToDomain()

	if err := h.projectService.Create(c.Request().Context(), project); err != nil {
		return err
	}

	return response.Created(c, dto.ProjectResFromDomain(project))
}

func (h *ProjectHandler) GetProject(c echo.Context) error {
	projectID, err := domain.ParseToUUIDProjectID(c.Param("id"))
	if err != nil {
		return err
	}

	project, err := h.projectService.FindByID(c.Request().Context(), projectID)
	if err != nil {
		return err
	}

	return response.OK(c, dto.ProjectResFromDomain(project))
}

func (h *ProjectHandler) UpdateProject(c echo.Context) error {
	projectID, err := domain.ParseToUUIDProjectID(c.Param("id"))
	if err != nil {
		return err
	}

	var req dto.UpdateProjectReq
	if err := reqctx.BindAndValidate(c, &req); err != nil {
		return err
	}

	// fetch existing so we can merge changes on top
	existing, err := h.projectService.FindByID(c.Request().Context(), projectID)
	if err != nil {
		return err
	}

	project := req.ToDomain(existing)
	project.ID = projectID

	if err := h.projectService.Update(c.Request().Context(), project); err != nil {
		return err
	}

	return response.OK(c, dto.ProjectResFromDomain(project))
}

func (h *ProjectHandler) DeleteProject(c echo.Context) error {
	projectID, err := domain.ParseToUUIDProjectID(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.projectService.Delete(c.Request().Context(), projectID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ProjectHandler) ListProjects(c echo.Context) error {
	pagination := reqctx.ParsePagination(c)
	filters := query.ParseFilters(c.QueryParams())
	sorts := query.ParseSort(c.QueryParam("sort"))

	projects, total, err := h.projectService.FindAll(
		c.Request().Context(),
		pagination.PerPage,
		pagination.Offset(),
		filters,
		sorts,
	)
	if err != nil {
		return err
	}

	return response.Paginated(c, dto.ProjectResFromDomainList(projects), pagination.Page, pagination.PerPage, total)
}
