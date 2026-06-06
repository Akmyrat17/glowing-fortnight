package persistance

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/skill/infra/persistance/dao"
	"github.com/boilerplate/internal/shared/pgutil"
	"github.com/boilerplate/pkg/query"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var skillAllowedFields = map[string]string{
	"id":         "id",
	"name":       "name",
	"level":      "level",
	"category":   "category",
	"years_exp":  "years_exp",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var skillColumns = []string{
	"id", "name", "level", "category",
	"icon_url", "years_exp",
	"created_at", "updated_at",
}

func scanSkill(row pgx.Row) (*dao.SkillDAO, error) {
	var s dao.SkillDAO
	err := row.Scan(
		&s.ID, &s.Name, &s.Level, &s.Category,
		&s.IconUrl, &s.YearsExp,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

type SkillRepoImpl struct {
	db *pgxpool.Pool
}

func NewSkillRepoImpl(db *pgxpool.Pool) *SkillRepoImpl {
	return &SkillRepoImpl{db: db}
}

func (r *SkillRepoImpl) Create(ctx context.Context, s *dao.SkillDAO) error {
	q, args, err := psql.Insert("skills").
		Columns(skillColumns...).
		Values(
			s.ID, s.Name, s.Level, s.Category,
			s.IconUrl, s.YearsExp,
			s.CreatedAt, s.UpdatedAt,
		).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, q, args...)
	return err
}

func (r *SkillRepoImpl) FindByID(ctx context.Context, id domain.SkillID) (*domain.Skill, error) {
	q, args, err := psql.Select(skillColumns...).
		From("skills").
		Where(sq.Eq{"id": pgutil.ToUUID(uuid.UUID(id))}).
		ToSql()
	if err != nil {
		return nil, err
	}
	s, err := scanSkill(r.db.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("skill with that id not found")
		}
		return nil, err
	}
	return s.ToDomain(), nil
}

func (r *SkillRepoImpl) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Skill, int64, error) {
	countBuilder := psql.Select("COUNT(*)").From("skills")
	countBuilder, err := query.ApplyFilters(countBuilder, filters, skillAllowedFields)
	if err != nil {
		return nil, 0, err
	}
	countQ, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := r.db.QueryRow(ctx, countQ, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	qb := psql.Select(skillColumns...).From("skills")
	qb, err = query.ApplyFilters(qb, filters, skillAllowedFields)
	if err != nil {
		return nil, 0, err
	}
	if len(sorts) > 0 {
		qb = query.ApplySort(qb, sorts, skillAllowedFields)
	} else {
		qb = qb.OrderBy("created_at DESC")
	}
	dbQ, args, err := qb.Limit(uint64(limit)).Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.db.Query(ctx, dbQ, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*domain.Skill
	for rows.Next() {
		s, err := scanSkill(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, s.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

func (r *SkillRepoImpl) FindAllGroupedByCategory(ctx context.Context) ([]*domain.Skill, error) {
	q, args, err := psql.Select(skillColumns...).From("skills").OrderBy("category", "name").ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*domain.Skill, 0)
	for rows.Next() {
		s, err := scanSkill(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, s.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *SkillRepoImpl) Update(ctx context.Context, s *dao.SkillDAO) error {
	setMap := sq.Eq{
		"name":       sq.Expr("COALESCE(NULLIF(?::text, ''), name)", s.Name),
		"level":      sq.Expr("COALESCE(NULLIF(?::text, ''), level)", s.Level),
		"category":   sq.Expr("COALESCE(NULLIF(?::text, ''), category)", s.Category),
		"icon_url":   s.IconUrl,
		"years_exp":  s.YearsExp,
		"updated_at": pgutil.ToTimestampt(time.Now()),
	}
	q, args, err := psql.Update("skills").
		SetMap(setMap).
		Where(sq.Eq{"id": s.ID}).
		ToSql()
	if err != nil {
		return err
	}
	result, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("skill with that id not found")
	}
	return nil
}

func (r *SkillRepoImpl) Delete(ctx context.Context, id domain.SkillID) (string, error) {
	// First fetch the icon_url
	var iconURL string
	fetchQ, fetchArgs, err := psql.Select("icon_url").
		From("skills").
		Where(sq.Eq{"id": pgutil.ToUUID(uuid.UUID(id))}).
		ToSql()
	if err != nil {
		return "", err
	}
	err = r.db.QueryRow(ctx, fetchQ, fetchArgs...).Scan(&iconURL)
	if err != nil {
		return "", fmt.Errorf("skill with that id not found")
	}

	// Then delete
	q, args, err := psql.Delete("skills").
		Where(sq.Eq{"id": pgutil.ToUUID(uuid.UUID(id))}).
		ToSql()
	if err != nil {
		return "", err
	}
	result, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return "", err
	}
	if result.RowsAffected() == 0 {
		return "", fmt.Errorf("skill with that id not found")
	}

	return iconURL, nil
}
