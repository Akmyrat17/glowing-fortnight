package domain

import (
	"time"

	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/google/uuid"
)

type ProjectType string
type ProjectStatus string

const (
	ProjectStatusDraft     ProjectStatus = "draft"
	ProjectStatusPublished ProjectStatus = "published"
	ProjectStatusArchived  ProjectStatus = "archived"
)

const (
	ProjectTypePet        ProjectType = "pet"
	ProjectTypeProduction ProjectType = "production"
)

type ProjectID uuid.UUID

type Project struct {
	ID          ProjectID
	Name        string
	Url         *string
	Description *string
	StartDate   *time.Time
	EndDate     *time.Time
	Tags        []string
	Status      ProjectStatus
	RepoUrl     *string
	ProjectType ProjectType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ParseToUUIDProjectID(s string) (ProjectID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return ProjectID{}, app_errors.ValidationError("invalid Project id")
	}
	return ProjectID(id), nil
}

func (p ProfileID) ParseToStringProjectID() string {
	return uuid.UUID(p).String()
}
func (s ProjectStatus) String() string { return string(s) }

func ParseProjectStatus(s string) (ProjectStatus, error) {
	switch ProjectStatus(s) {
	case ProjectStatusDraft, ProjectStatusPublished, ProjectStatusArchived:
		return ProjectStatus(s), nil
	}
	return "", app_errors.ValidationError("invalid project status, must be one of: draft, published, archived")
}

func (t ProjectType) String() string { return string(t) }

func ParseProjectType(s string) (ProjectType, error) {
	switch ProjectType(s) {
	case ProjectTypePet, ProjectTypeProduction:
		return ProjectType(s), nil
	}
	return "", app_errors.ValidationError("invalid project type, must be one of: pet, production")
}

func NewProject(name string, url, description, repoUrl *string, startDate, endDate *time.Time, tags []string, status ProjectStatus, projectType ProjectType) *Project {
	now := time.Now()
	return &Project{
		ID:          ProjectID(uuid.New()),
		Name:        name,
		Description: description,
		RepoUrl:     repoUrl,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      status,
		ProjectType: projectType,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p *Project) MarkUpdated() {
	p.UpdatedAt = time.Now()
}
