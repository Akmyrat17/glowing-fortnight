package application

import (
	"context"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/experience/infra/persistance/dao"
	"github.com/boilerplate/pkg/query"
)

type ExperienceRepository interface {
	Create(ctx context.Context, e *dao.ExperienceDAO) error
	FindByID(ctx context.Context, id domain.ExperienceID) (*domain.Experience, error)
	FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Experience, int64, error)
	Update(ctx context.Context, e *dao.ExperienceDAO) error
	Delete(ctx context.Context, id domain.ExperienceID) error
}

type ExperienceService struct {
	repo ExperienceRepository
}

func NewExperienceService(repo ExperienceRepository) *ExperienceService {
	return &ExperienceService{repo: repo}
}

func (s *ExperienceService) Create(ctx context.Context, e *domain.Experience) error {
	return s.repo.Create(ctx, dao.ExperienceFromDomain(e))
}

func (s *ExperienceService) FindByID(ctx context.Context, id domain.ExperienceID) (*domain.Experience, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ExperienceService) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Experience, int64, error) {
	return s.repo.FindAll(ctx, limit, offset, filters, sorts)
}

func (s *ExperienceService) Update(ctx context.Context, e *domain.Experience) error {
	return s.repo.Update(ctx, dao.ExperienceFromDomain(e))
}

func (s *ExperienceService) Delete(ctx context.Context, id domain.ExperienceID) error {
	return s.repo.Delete(ctx, id)
}
