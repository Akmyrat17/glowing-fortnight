package dao

import (
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type EducationDAO struct {
	ID           pgtype.UUID
	School       string
	Degree       pgtype.Text
	FieldOfStudy pgtype.Text
	StartDate    pgtype.Timestamptz
	EndDate      pgtype.Timestamptz
	IsCurrent    bool
	Description  pgtype.Text
	CreatedAt    pgtype.Timestamptz
	UpdatedAt    pgtype.Timestamptz
}

func (d *EducationDAO) ToDomain() *domain.Education {
	return &domain.Education{
		ID:           domain.EducationID(pgutil.FromUUID(d.ID)),
		School:       d.School,
		Degree:       pgutil.FromNullableText(d.Degree),
		FieldOfStudy: pgutil.FromNullableText(d.FieldOfStudy),
		StartDate:    pgutil.FromTimestampt(d.StartDate),
		EndDate:      pgutil.FromNullableTimestampt(d.EndDate),
		IsCurrent:    d.IsCurrent,
		Description:  pgutil.FromNullableText(d.Description),
		CreatedAt:    pgutil.FromTimestampt(d.CreatedAt),
		UpdatedAt:    pgutil.FromTimestampt(d.UpdatedAt),
	}
}

func EducationFromDomain(e *domain.Education) *EducationDAO {
	return &EducationDAO{
		ID:           pgutil.ToUUID(uuid.UUID(e.ID)),
		School:       e.School,
		Degree:       pgutil.ToNullableText(e.Degree),
		FieldOfStudy: pgutil.ToNullableText(e.FieldOfStudy),
		StartDate:    pgutil.ToTimestampt(e.StartDate),
		EndDate:      pgutil.ToNullableTimestampt(e.EndDate),
		IsCurrent:    e.IsCurrent,
		Description:  pgutil.ToNullableText(e.Description),
		CreatedAt:    pgutil.ToTimestampt(e.CreatedAt),
		UpdatedAt:    pgutil.ToTimestampt(e.UpdatedAt),
	}
}
