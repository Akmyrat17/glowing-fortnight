package persistance

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/skill/infra/persistance/dao"
	"github.com/boilerplate/internal/shared/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func scanSkill(row pgx.Row) (*dao.SkillDAO, error) {
	var skill dao.SkillDAO
	err := row.Scan(
		&skill.ID, &skill.Name, &skill.Level, &skill.Category,
		&skill.IconUrl, &skill.YearsExp,
		&skill.CreatedAt, &skill.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &skill, nil
}
func (r *ProjectRepoImpl) AddSkills(ctx context.Context, projectID domain.ProjectID, skillIDs []domain.SkillID) error {
	if len(skillIDs) == 0 {
		return nil
	}

	builder := psql.Insert("project_skills").Columns("project_id", "skill_id")
	for _, skillID := range skillIDs {
		builder = builder.Values(
			pgutil.ToUUID(uuid.UUID(projectID)),
			pgutil.ToUUID(uuid.UUID(skillID)),
		)
	}
	q, args, err := builder.Suffix("ON CONFLICT DO NOTHING").ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, q, args...)
	return err
}

func (r *ProjectRepoImpl) RemoveSkills(ctx context.Context, projectID domain.ProjectID, skillIDs []domain.SkillID) error {
	if len(skillIDs) == 0 {
		return nil
	}

	ids := make([]uuid.UUID, len(skillIDs))
	for i, skillID := range skillIDs {
		ids[i] = uuid.UUID(skillID)
	}

	q, args, err := psql.Delete("project_skills").
		Where(sq.Eq{
			"project_id": pgutil.ToUUID(uuid.UUID(projectID)),
			"skill_id":   ids,
		}).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, q, args...)
	return err
}

func (r *ProjectRepoImpl) FindSkills(ctx context.Context, projectID domain.ProjectID) ([]*domain.Skill, error) {
	q, args, err := psql.Select(
		"s.id", "s.name", "s.level", "s.category",
		"s.icon_url", "s.years_exp", "s.created_at", "s.updated_at",
	).
		From("skills s").
		Join("project_skills ps ON ps.skill_id = s.id").
		Where(sq.Eq{"ps.project_id": pgutil.ToUUID(uuid.UUID(projectID))}).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var skills []*domain.Skill
	for rows.Next() {
		s, err := scanSkill(rows)
		if err != nil {
			return nil, err
		}
		skills = append(skills, s.ToDomain())
	}
	return skills, rows.Err()
}
