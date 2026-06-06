package domain

import (
	"time"

	"github.com/boilerplate/internal/shared/app_errors"
	"github.com/google/uuid"
)

type SkillID uuid.UUID

func (s SkillID) String() string { return uuid.UUID(s).String() }
func ParseSkillID(str string) (SkillID, error) {
	id, err := uuid.Parse(str)
	if err != nil {
		return SkillID{}, app_errors.ValidationError("invalid skill ID format")
	}
	return SkillID(id), nil
}

type SkillLevel string

const (
	SkillLevelBeginner     SkillLevel = "beginner"
	SkillLevelIntermediate SkillLevel = "intermediate"
	SkillLevelAdvanced     SkillLevel = "advanced"
	SkillLevelExpert       SkillLevel = "expert"
)

func (l SkillLevel) Valid() bool {
	switch l {
	case SkillLevelBeginner, SkillLevelIntermediate, SkillLevelAdvanced, SkillLevelExpert:
		return true
	}
	return false
}

type SkillCategory string

const (
	SkillCategoryLanguage  SkillCategory = "language"
	SkillCategoryFramework SkillCategory = "framework"
	SkillCategoryDatabase  SkillCategory = "database"
	SkillCategoryTool      SkillCategory = "tool"
	SkillCategoryConcept   SkillCategory = "concept"
	SkillCategoryORM       SkillCategory = "orm"
)

func (c SkillCategory) Valid() bool {
	switch c {
	case SkillCategoryLanguage, SkillCategoryFramework, SkillCategoryDatabase, SkillCategoryTool, SkillCategoryConcept, SkillCategoryORM:
		return true
	}
	return false
}

type Skill struct {
	ID        SkillID
	Name      string
	Level     SkillLevel
	Category  SkillCategory
	IconUrl   *string
	YearsExp  *int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSkill(name string, level SkillLevel, category SkillCategory, iconUrl *string, yearsExp *int) (*Skill, error) {
	if !level.Valid() {
		return nil, app_errors.ValidationError("invalid skill level")
	}
	if !category.Valid() {
		return nil, app_errors.ValidationError("invalid skill category")
	}
	return &Skill{
		ID:        SkillID(uuid.New()),
		Name:      name,
		Level:     level,
		Category:  category,
		IconUrl:   iconUrl,
		YearsExp:  yearsExp,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *Skill) MarkUpdated() {
	s.UpdatedAt = time.Now()
}
