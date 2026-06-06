package dao

import (
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SkillDAO struct {
	ID        pgtype.UUID
	Name      string
	Level     string
	Category  string
	IconUrl   pgtype.Text
	YearsExp  pgtype.Int4
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func (d *SkillDAO) ToDomain() *domain.Skill {
	var yearsExp *int
	if d.YearsExp.Valid {
		v := int(d.YearsExp.Int32)
		yearsExp = &v
	}
	return &domain.Skill{
		ID:        domain.SkillID(pgutil.FromUUID(d.ID)),
		Name:      d.Name,
		Level:     domain.SkillLevel(d.Level),
		Category:  domain.SkillCategory(d.Category),
		IconUrl:   pgutil.FromNullableText(d.IconUrl),
		YearsExp:  yearsExp,
		CreatedAt: pgutil.FromTimestampt(d.CreatedAt),
		UpdatedAt: pgutil.FromTimestampt(d.UpdatedAt),
	}
}

func SkillFromDomain(s *domain.Skill) *SkillDAO {
	var yearsExp pgtype.Int4
	if s.YearsExp != nil {
		yearsExp = pgtype.Int4{Int32: int32(*s.YearsExp), Valid: true}
	}
	return &SkillDAO{
		ID:        pgutil.ToUUID(uuid.UUID(s.ID)),
		Name:      s.Name,
		Level:     string(s.Level),
		Category:  string(s.Category),
		IconUrl:   pgutil.ToNullableText(s.IconUrl),
		YearsExp:  yearsExp,
		CreatedAt: pgutil.ToTimestampt(s.CreatedAt),
		UpdatedAt: pgutil.ToTimestampt(s.UpdatedAt),
	}
}
