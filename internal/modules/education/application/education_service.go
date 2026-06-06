package application

import (
	"context"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/education/infra/persistance/dao"
	"github.com/boilerplate/pkg/query"
)

type EducationRepository interface {
	Create(ctx context.Context, e *dao.EducationDAO) error
	FindByID(ctx context.Context, id domain.EducationID) (*domain.Education, error)
	FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Education, int64, error)
	Update(ctx context.Context, e *dao.EducationDAO) error
	Delete(ctx context.Context, id domain.EducationID) error
}

type EducationService struct {
	repo EducationRepository
}

func NewEducationService(repo EducationRepository) *EducationService {
	return &EducationService{repo: repo}
}

func (s *EducationService) Create(ctx context.Context, e *domain.Education) error {
	return s.repo.Create(ctx, dao.EducationFromDomain(e))
}

func (s *EducationService) FindByID(ctx context.Context, id domain.EducationID) (*domain.Education, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *EducationService) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Education, int64, error) {
	return s.repo.FindAll(ctx, limit, offset, filters, sorts)
}

func (s *EducationService) Update(ctx context.Context, e *domain.Education) error {
	return s.repo.Update(ctx, dao.EducationFromDomain(e))
}

func (s *EducationService) Delete(ctx context.Context, id domain.EducationID) error {
	return s.repo.Delete(ctx, id)
}
