package dto

import (
	"time"

	"github.com/boilerplate/internal/domain"
	"github.com/google/uuid"
)

// ---------------------------------------------------------------
// REQUEST DTOs
// ---------------------------------------------------------------

type CreateExperienceReq struct {
	Company     string     `json:"company"   validate:"required"`
	Position    string     `json:"position"  validate:"required"`
	Description *string    `json:"description"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date"`
	Location    *string    `json:"location"`
	ProjectID   *uuid.UUID `json:"project_id"`
}

type UpdateExperienceReq struct {
	Company     string     `json:"company"`
	Position    string     `json:"position"`
	Description *string    `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Location    *string    `json:"location"`
	ProjectID   *uuid.UUID `json:"project_id"`
}

// ---------------------------------------------------------------
// REQUEST → DOMAIN
// ---------------------------------------------------------------

func (r *CreateExperienceReq) ToDomain() *domain.Experience {
	return domain.NewExperience(
		r.Company,
		r.Position,
		r.Description,
		r.Location,
		r.StartDate,
		r.EndDate,
		r.ProjectID,
	)
}

func (r *UpdateExperienceReq) ToDomain(existing *domain.Experience) *domain.Experience {
	if r.Company != "" {
		existing.Company = r.Company
	}
	if r.Position != "" {
		existing.Position = r.Position
	}
	if r.Description != nil {
		existing.Description = r.Description
	}
	if r.StartDate != nil {
		existing.StartDate = *r.StartDate
	}
	if r.EndDate != nil {
		existing.EndDate = r.EndDate
		existing.IsCurrent = false
	}
	if r.Location != nil {
		existing.Location = r.Location
	}

	existing.MarkUpdated()
	return existing
}

// ---------------------------------------------------------------
// RESPONSE DTO
// ---------------------------------------------------------------

type ExperienceRes struct {
	ID          uuid.UUID           `json:"id"`
	Company     string              `json:"company"`
	Position    string              `json:"position"`
	Description *string             `json:"description"`
	StartDate   time.Time           `json:"start_date"`
	EndDate     *time.Time          `json:"end_date"`
	IsCurrent   bool                `json:"is_current"`
	Location    *string             `json:"location"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	ProjectID   *uuid.UUID          `json:"project_id"`
	Project     *EmbeddedProjectRes `json:"project,omitempty"` // omitempty — not always loaded
}

// ---------------------------------------------------------------
// EMBEDDED PROJECT (used inside experience response only)
// ---------------------------------------------------------------

type EmbeddedProjectRes struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Url         *string   `json:"url"`
	RepoUrl     *string   `json:"repo_url"`
	Status      string    `json:"status"`
	ProjectType string    `json:"project_type"`
}

func embeddedProjectFromDomain(p *domain.Project) *EmbeddedProjectRes {
	return &EmbeddedProjectRes{
		ID:          uuid.UUID(p.ID),
		Name:        p.Name,
		Url:         p.Url,
		RepoUrl:     p.RepoUrl,
		Status:      string(p.Status),
		ProjectType: string(p.ProjectType),
	}
}

// ---------------------------------------------------------------
// DOMAIN → RESPONSE
// ---------------------------------------------------------------
func ExperienceResFromDomain(e *domain.Experience) *ExperienceRes {
	res := &ExperienceRes{
		ID:          uuid.UUID(e.ID),
		Company:     e.Company,
		Position:    e.Position,
		Description: e.Description,
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
		IsCurrent:   e.IsCurrent,
		Location:    e.Location,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}

	if e.ProjectID != nil {
		id := uuid.UUID(*e.ProjectID)
		res.ProjectID = &id
	}
	if e.Project != nil {
		res.Project = embeddedProjectFromDomain(e.Project)
	}

	return res
}

func ExperienceResFromDomainList(list []*domain.Experience) []*ExperienceRes {
	res := make([]*ExperienceRes, len(list))
	for i, e := range list {
		res[i] = ExperienceResFromDomain(e)
	}
	return res
}
