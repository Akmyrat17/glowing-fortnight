package dto

import (
	"time"

	"github.com/boilerplate/internal/domain"
	"github.com/google/uuid"
)

// ---------------------------------------------------------------
// REQUEST DTOs  — what the handler receives from the client
// ---------------------------------------------------------------

type CreateProjectReq struct {
	Name        string     `json:"name"         validate:"required"`
	Description *string    `json:"description"`
	Url         *string    `json:"url"`
	RepoUrl     *string    `json:"repo_url"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Tags        []string   `json:"tags"`
	Status      string     `json:"status"       validate:"required"`
	ProjectType string     `json:"project_type" validate:"required"`
}

type UpdateProjectReq struct {
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Url         *string    `json:"url"`
	RepoUrl     *string    `json:"repo_url"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Tags        []string   `json:"tags"`
	Status      string     `json:"status"`
	ProjectType string     `json:"project_type"`
}

// ---------------------------------------------------------------
// REQUEST → DOMAIN transforms
// ---------------------------------------------------------------

func (r *CreateProjectReq) ToDomain() *domain.Project {
	status := domain.ProjectStatus(r.Status)
	projectType := domain.ProjectType(r.ProjectType)

	return domain.NewProject(
		r.Name,
		r.Url,
		r.Description,
		r.RepoUrl,
		r.StartDate,
		r.EndDate,
		r.Tags,
		status,
		projectType,
	)
}

// existing is the current DB state — we merge changes on top of it
func (r *UpdateProjectReq) ToDomain(existing *domain.Project) *domain.Project {
	if r.Name != "" {
		existing.Name = r.Name
	}
	if r.Description != nil {
		existing.Description = r.Description
	}
	if r.Url != nil {
		existing.Url = r.Url
	}
	if r.RepoUrl != nil {
		existing.RepoUrl = r.RepoUrl
	}
	if r.StartDate != nil {
		existing.StartDate = r.StartDate
	}
	if r.EndDate != nil {
		existing.EndDate = r.EndDate
	}
	if r.Tags != nil {
		existing.Tags = r.Tags
	}
	if r.Status != "" {
		existing.Status = domain.ProjectStatus(r.Status)
	}
	if r.ProjectType != "" {
		existing.ProjectType = domain.ProjectType(r.ProjectType)
	}

	existing.MarkUpdated()
	return existing
}

// ---------------------------------------------------------------
// RESPONSE DTO  — what we send back to the client
// ---------------------------------------------------------------

type ProjectRes struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Url         *string    `json:"url"`
	RepoUrl     *string    `json:"repo_url"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Tags        []string   `json:"tags"`
	Status      string     `json:"status"`
	ProjectType string     `json:"project_type"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ---------------------------------------------------------------
// DOMAIN → RESPONSE transforms
// ---------------------------------------------------------------

func ProjectResFromDomain(p *domain.Project) *ProjectRes {
	return &ProjectRes{
		ID:          uuid.UUID(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Url:         p.Url,
		RepoUrl:     p.RepoUrl,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		Tags:        p.Tags,
		Status:      string(p.Status),
		ProjectType: string(p.ProjectType),
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ProjectResFromDomainList(projects []*domain.Project) []*ProjectRes {
	res := make([]*ProjectRes, len(projects))
	for i, p := range projects {
		res[i] = ProjectResFromDomain(p)
	}
	return res
}

type ProjectSkillsReq struct {
	SkillIDs []uuid.UUID `json:"skill_ids" validate:"required,min=1,dive,required"`
}
