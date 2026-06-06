package application

import (
	"context"
	"os"
	"path/filepath"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/skill/infra/persistance/dao"
	"github.com/boilerplate/pkg/query"
)

type SkillRepository interface {
	Create(ctx context.Context, s *dao.SkillDAO) error
	FindByID(ctx context.Context, id domain.SkillID) (*domain.Skill, error)
	FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Skill, int64, error)
	FindAllGroupedByCategory(ctx context.Context) ([]*domain.Skill, error)
	Update(ctx context.Context, s *dao.SkillDAO) error
	Delete(ctx context.Context, id domain.SkillID) (string, error)
}

type SkillService struct {
	repo SkillRepository
}

func NewSkillService(repo SkillRepository) *SkillService {
	return &SkillService{repo: repo}
}

func (s *SkillService) Create(ctx context.Context, skill *domain.Skill) error {
	return s.repo.Create(ctx, dao.SkillFromDomain(skill))
}

func (s *SkillService) FindByID(ctx context.Context, id domain.SkillID) (*domain.Skill, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *SkillService) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Skill, int64, error) {
	return s.repo.FindAll(ctx, limit, offset, filters, sorts)
}

func (s *SkillService) FindAllGroupedByCategory(ctx context.Context) (map[domain.SkillCategory][]*domain.Skill, error) {
	skills, err := s.repo.FindAllGroupedByCategory(ctx)
	if err != nil {
		return nil, err
	}

	groups := make(map[domain.SkillCategory][]*domain.Skill)
	for _, skill := range skills {
		groups[skill.Category] = append(groups[skill.Category], skill)
	}
	return groups, nil
}
func (s *SkillService) Update(ctx context.Context, skill *domain.Skill) error {
	return s.repo.Update(ctx, dao.SkillFromDomain(skill))
}

func (s *SkillService) Delete(ctx context.Context, id domain.SkillID) error {
	iconURL, err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Remove the file from disk (non-blocking, ignore error if file missing)
	if iconURL != "" {
		filePath := filepath.Join(".", iconURL) // e.g. "./uploads/skill/xxx.jpg"
		_ = os.Remove(filePath)
	}

	return nil
}
