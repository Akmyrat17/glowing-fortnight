package dto

import (
	"time"

	"github.com/boilerplate/internal/domain"
	"github.com/google/uuid"
)

// ---------------------------------------------------------------
// REQUEST DTOs
// ---------------------------------------------------------------

type CreateEducationReq struct {
	School       string     `json:"school"         validate:"required"`
	Degree       *string    `json:"degree"`
	FieldOfStudy *string    `json:"field_of_study"`
	StartDate    time.Time  `json:"start_date"     validate:"required"`
	EndDate      *time.Time `json:"end_date"`
	Description  *string    `json:"description"`
}

type UpdateEducationReq struct {
	School       string     `json:"school"`
	Degree       *string    `json:"degree"`
	FieldOfStudy *string    `json:"field_of_study"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Description  *string    `json:"description"`
}

// ---------------------------------------------------------------
// REQUEST → DOMAIN
// ---------------------------------------------------------------

func (r *CreateEducationReq) ToDomain() *domain.Education {
	return domain.NewEducation(
		r.School,
		r.Degree,
		r.FieldOfStudy,
		r.Description,
		r.StartDate,
		r.EndDate,
	)
}

func (r *UpdateEducationReq) ToDomain(existing *domain.Education) *domain.Education {
	if r.School != "" {
		existing.School = r.School
	}
	if r.Degree != nil {
		existing.Degree = r.Degree
	}
	if r.FieldOfStudy != nil {
		existing.FieldOfStudy = r.FieldOfStudy
	}
	if r.StartDate != nil {
		existing.StartDate = *r.StartDate
	}
	if r.EndDate != nil {
		existing.EndDate = r.EndDate
		existing.IsCurrent = false
	}
	if r.Description != nil {
		existing.Description = r.Description
	}

	existing.MarkUpdated()
	return existing
}

// ---------------------------------------------------------------
// RESPONSE DTO
// ---------------------------------------------------------------

type EducationRes struct {
	ID           uuid.UUID  `json:"id"`
	School       string     `json:"school"`
	Degree       *string    `json:"degree"`
	FieldOfStudy *string    `json:"field_of_study"`
	StartDate    time.Time  `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	IsCurrent    bool       `json:"is_current"`
	Description  *string    `json:"description"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// ---------------------------------------------------------------
// DOMAIN → RESPONSE
// ---------------------------------------------------------------

func EducationResFromDomain(e *domain.Education) *EducationRes {
	return &EducationRes{
		ID:           uuid.UUID(e.ID),
		School:       e.School,
		Degree:       e.Degree,
		FieldOfStudy: e.FieldOfStudy,
		StartDate:    e.StartDate,
		EndDate:      e.EndDate,
		IsCurrent:    e.IsCurrent,
		Description:  e.Description,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func EducationResFromDomainList(list []*domain.Education) []*EducationRes {
	res := make([]*EducationRes, len(list))
	for i, e := range list {
		res[i] = EducationResFromDomain(e)
	}
	return res
}
