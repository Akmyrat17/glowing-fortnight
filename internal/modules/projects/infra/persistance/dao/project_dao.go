package dao

import (
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectDAO struct {
	ID          pgtype.UUID
	Name        string
	Description pgtype.Text
	Url         pgtype.Text
	RepoUrl     pgtype.Text
	StartDate   pgtype.Timestamptz
	EndDate     pgtype.Timestamptz
	Tags        pgtype.Array[string]
	Status      string
	ProjectType string
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
}

func (d *ProjectDAO) ToDomain() *domain.Project {
	return &domain.Project{
		ID:          domain.ProjectID(pgutil.FromUUID(d.ID)),
		Name:        d.Name,
		Description: pgutil.FromNullableText(d.Description),
		Url:         pgutil.FromNullableText(d.Url),
		RepoUrl:     pgutil.FromNullableText(d.RepoUrl),
		StartDate:   pgutil.FromNullableTimestampt(d.StartDate),
		EndDate:     pgutil.FromNullableTimestampt(d.EndDate),
		Tags:        pgutil.FromPostgresArray(d.Tags),
		Status:      domain.ProjectStatus(d.Status),
		ProjectType: domain.ProjectType(d.ProjectType),
		CreatedAt:   pgutil.FromTimestampt(d.CreatedAt),
		UpdatedAt:   pgutil.FromTimestampt(d.UpdatedAt),
	}
}

func FromDomain(project *domain.Project) *ProjectDAO {
	return &ProjectDAO{
		ID:          pgutil.ToUUID(uuid.UUID(project.ID)),
		Name:        project.Name,
		Description: pgutil.ToNullableText(project.Description),
		Url:         pgutil.ToNullableText(project.Url),
		RepoUrl:     pgutil.ToNullableText(project.RepoUrl),
		StartDate:   pgutil.ToNullableTimestampt(project.StartDate),
		EndDate:     pgutil.ToNullableTimestampt(project.EndDate),
		Tags:        pgutil.ToPostgresArray(project.Tags),
		Status:      string(project.Status),
		ProjectType: string(project.ProjectType),
		CreatedAt:   pgutil.ToTimestampt(project.CreatedAt),
		UpdatedAt:   pgutil.ToTimestampt(project.UpdatedAt),
	}
}
