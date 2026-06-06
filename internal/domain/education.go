package domain

import (
	"time"

	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/google/uuid"
)

type EducationID uuid.UUID

func (e EducationID) String() string { return uuid.UUID(e).String() }
func ParseEducationID(s string) (EducationID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return EducationID{}, app_errors.ValidationError("invalid education ID format")
	}
	return EducationID(id), nil
}

type Education struct {
	ID           EducationID
	School       string
	Degree       *string
	FieldOfStudy *string
	StartDate    time.Time
	EndDate      *time.Time
	IsCurrent    bool
	Description  *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewEducation(school string, degree, fieldOfStudy, description *string, startDate time.Time, endDate *time.Time) *Education {
	return &Education{
		ID:           EducationID(uuid.New()),
		School:       school,
		Degree:       degree,
		FieldOfStudy: fieldOfStudy,
		Description:  description,
		StartDate:    startDate,
		EndDate:      endDate,
		IsCurrent:    endDate == nil,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (e *Education) MarkUpdated() {
	e.UpdatedAt = time.Now()
}
