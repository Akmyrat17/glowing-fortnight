package application

import (
	"context"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/projects/infra/persistance/dao"
	"github.com/boilerplate/pkg/query"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *dao.ProjectDAO) error
	FindByID(ctx context.Context, id domain.ProjectID) (*domain.Project, error)
	FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Project, int64, error)
	Update(ctx context.Context, project *dao.ProjectDAO) error
	Delete(ctx context.Context, id domain.ProjectID) error
}

type ProjectService struct {
	repo ProjectRepository
}

func NewProjectService(repo ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, project *domain.Project) error {
	return s.repo.Create(ctx, dao.FromDomain(project))
}

func (s *ProjectService) FindByID(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ProjectService) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Project, int64, error) {
	return s.repo.FindAll(ctx, limit, offset, filters, sorts)
}

func (s *ProjectService) Update(ctx context.Context, project *domain.Project) error {
	return s.repo.Update(ctx, dao.FromDomain(project))
}

func (s *ProjectService) Delete(ctx context.Context, id domain.ProjectID) error {
	return s.repo.Delete(ctx, id)
}
