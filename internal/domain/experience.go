package domain

import (
	"time"

	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/google/uuid"
)

type ExperienceID uuid.UUID

func (e ExperienceID) String() string { return uuid.UUID(e).String() }
func ParseExperienceID(s string) (ExperienceID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ExperienceID{}, app_errors.ValidationError("invalid experience ID format")
	}
	return ExperienceID(id), nil
}

type Experience struct {
	ID          ExperienceID
	Company     string
	Position    string
	Description *string
	StartDate   time.Time
	EndDate     *time.Time
	IsCurrent   bool
	Location    *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ProjectID   *ProjectID
	Project     *Project
}

func NewExperience(company, position string, description, location *string, startDate time.Time, endDate *time.Time, projectID *uuid.UUID) *Experience {
	var pid *ProjectID
	if projectID != nil {
		id := ProjectID(*projectID) // dereference uuid.UUID, cast to ProjectID
		pid = &id
	}
	return &Experience{
		ID:          ExperienceID(uuid.New()),
		Company:     company,
		Position:    position,
		Description: description,
		Location:    location,
		StartDate:   startDate,
		EndDate:     endDate,
		IsCurrent:   endDate == nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ProjectID:   pid,
	}
}

func (e *Experience) MarkUpdated() {
	e.UpdatedAt = time.Now()
}
