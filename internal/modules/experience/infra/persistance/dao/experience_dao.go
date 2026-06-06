package dao

import (
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ExperienceDAO struct {
	ID          pgtype.UUID
	Company     string
	Position    string
	Description pgtype.Text
	StartDate   pgtype.Timestamptz
	EndDate     pgtype.Timestamptz
	IsCurrent   bool
	Location    pgtype.Text
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	ProjectID   pgtype.UUID
}

func (d *ExperienceDAO) ToDomain() *domain.Experience {
	rawProjectId := pgutil.FromNullableUUID(d.ProjectID)
	var projectID *domain.ProjectID
	if rawProjectId != nil {
		id := domain.ProjectID(*rawProjectId)
		projectID = &id
	}
	return &domain.Experience{
		ID:          domain.ExperienceID(pgutil.FromUUID(d.ID)),
		Company:     d.Company,
		Position:    d.Position,
		Description: pgutil.FromNullableText(d.Description),
		StartDate:   pgutil.FromTimestampt(d.StartDate),
		EndDate:     pgutil.FromNullableTimestampt(d.EndDate),
		IsCurrent:   d.IsCurrent,
		Location:    pgutil.FromNullableText(d.Location),
		CreatedAt:   pgutil.FromTimestampt(d.CreatedAt),
		UpdatedAt:   pgutil.FromTimestampt(d.UpdatedAt),
		ProjectID:   projectID,
	}
}

func ExperienceFromDomain(e *domain.Experience) *ExperienceDAO {
	return &ExperienceDAO{
		ID:          pgutil.ToUUID(uuid.UUID(e.ID)),
		Company:     e.Company,
		Position:    e.Position,
		Description: pgutil.ToNullableText(e.Description),
		StartDate:   pgutil.ToTimestampt(e.StartDate),
		EndDate:     pgutil.ToNullableTimestampt(e.EndDate),
		IsCurrent:   e.IsCurrent,
		Location:    pgutil.ToNullableText(e.Location),
		CreatedAt:   pgutil.ToTimestampt(e.CreatedAt),
		UpdatedAt:   pgutil.ToTimestampt(e.UpdatedAt),
		ProjectID:   pgutil.ToNullableUUID((*uuid.UUID)(e.ProjectID)),
	}
}
