package dto

import (
	"time"

	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/google/uuid"
)

// ---------------------------------------------------------------
// REQUEST DTOs
// ---------------------------------------------------------------

type CreateSkillReq struct {
	Name     string  `json:"name"     validate:"required"`
	Level    string  `json:"level"    validate:"required"`
	Category string  `json:"category" validate:"required"`
	IconUrl  *string `json:"icon_url"`
	YearsExp *int    `json:"years_exp"`
}

type UpdateSkillReq struct {
	Name     string  `json:"name"`
	Level    string  `json:"level"`
	Category string  `json:"category"`
	IconUrl  *string `json:"icon_url"`
	YearsExp *int    `json:"years_exp"`
}

// ---------------------------------------------------------------
// REQUEST → DOMAIN
// ---------------------------------------------------------------

func (r *CreateSkillReq) ToDomain() (*domain.Skill, error) {
	return domain.NewSkill(
		r.Name,
		domain.SkillLevel(r.Level),
		domain.SkillCategory(r.Category),
		r.IconUrl,
		r.YearsExp,
	)
}

func (r *UpdateSkillReq) ToDomain(existing *domain.Skill) (*domain.Skill, error) {
	if r.Name != "" {
		existing.Name = r.Name
	}
	if r.Level != "" {
		level := domain.SkillLevel(r.Level)
		if !level.Valid() {
			return nil, app_errors.ValidationError("invalid skill level")
		}
		existing.Level = level
	}
	if r.Category != "" {
		category := domain.SkillCategory(r.Category)
		if !category.Valid() {
			return nil, app_errors.ValidationError("invalid skill category")
		}
		existing.Category = category
	}
	if r.IconUrl != nil {
		existing.IconUrl = r.IconUrl
	}
	if r.YearsExp != nil {
		existing.YearsExp = r.YearsExp
	}

	existing.MarkUpdated()
	return existing, nil
}

// ---------------------------------------------------------------
// RESPONSE DTO
// ---------------------------------------------------------------

type SkillRes struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Level     string    `json:"level"`
	Category  string    `json:"category"`
	IconUrl   *string   `json:"icon_url"`
	YearsExp  *int      `json:"years_exp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ---------------------------------------------------------------
// DOMAIN → RESPONSE
// ---------------------------------------------------------------

func SkillResFromDomain(s *domain.Skill) *SkillRes {
	return &SkillRes{
		ID:        uuid.UUID(s.ID),
		Name:      s.Name,
		Level:     string(s.Level),
		Category:  string(s.Category),
		IconUrl:   s.IconUrl,
		YearsExp:  s.YearsExp,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func SkillResFromDomainList(list []*domain.Skill) []*SkillRes {
	res := make([]*SkillRes, len(list))
	for i, s := range list {
		res[i] = SkillResFromDomain(s)
	}
	return res
}

type SkillCategoryGroupRes struct {
	Category string      `json:"category"`
	Skills   []*SkillRes `json:"skills"`
}
