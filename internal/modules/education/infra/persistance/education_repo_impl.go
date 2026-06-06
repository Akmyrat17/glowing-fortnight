package persistance

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/education/infra/persistance/dao"
	"github.com/boilerplate/internal/shared/pgutil"
	"github.com/boilerplate/pkg/query"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var educationAllowedFields = map[string]string{
	"id":             "id",
	"school":         "school",
	"degree":         "degree",
	"field_of_study": "field_of_study",
	"is_current":     "is_current",
	"start_date":     "start_date",
	"end_date":       "end_date",
	"created_at":     "created_at",
	"updated_at":     "updated_at",
}

var educationColumns = []string{
	"id", "school", "degree", "field_of_study",
	"start_date", "end_date", "is_current", "description",
	"created_at", "updated_at",
}

func scanEducation(row pgx.Row) (*dao.EducationDAO, error) {
	var e dao.EducationDAO
	err := row.Scan(
		&e.ID, &e.School, &e.Degree, &e.FieldOfStudy,
		&e.StartDate, &e.EndDate, &e.IsCurrent, &e.Description,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

type EducationRepoImpl struct {
	db *pgxpool.Pool
}

func NewEducationRepoImpl(db *pgxpool.Pool) *EducationRepoImpl {
	return &EducationRepoImpl{db: db}
}

func (r *EducationRepoImpl) Create(ctx context.Context, e *dao.EducationDAO) error {
	q, args, err := psql.Insert("educations").
		Columns(educationColumns...).
		Values(
			e.ID, e.School, e.Degree, e.FieldOfStudy,
			e.StartDate, e.EndDate, e.IsCurrent, e.Description,
			e.CreatedAt, e.UpdatedAt,
		).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, q, args...)
	return err
}

func (r *EducationRepoImpl) FindByID(ctx context.Context, id domain.EducationID) (*domain.Education, error) {
	q, args, err := psql.Select(educationColumns...).
		From("educations").
		Where(sq.Eq{"id": pgutil.ToUUID(uuid.UUID(id))}).
		ToSql()
	if err != nil {
		return nil, err
	}
	e, err := scanEducation(r.db.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("education with that id not found")
		}
		return nil, err
	}
	return e.ToDomain(), nil
}

func (r *EducationRepoImpl) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Education, int64, error) {
	countBuilder := psql.Select("COUNT(*)").From("educations")
	countBuilder, err := query.ApplyFilters(countBuilder, filters, educationAllowedFields)
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

	qb := psql.Select(educationColumns...).From("educations")
	qb, err = query.ApplyFilters(qb, filters, educationAllowedFields)
	if err != nil {
		return nil, 0, err
	}
	if len(sorts) > 0 {
		qb = query.ApplySort(qb, sorts, educationAllowedFields)
	} else {
		qb = qb.OrderBy("start_date DESC")
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

	var result []*domain.Education
	for rows.Next() {
		e, err := scanEducation(rows)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, e.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return result, total, nil
}

func (r *EducationRepoImpl) Update(ctx context.Context, e *dao.EducationDAO) error {
	setMap := sq.Eq{
		"school":         sq.Expr("COALESCE(NULLIF(?::text, ''), school)", e.School),
		"degree":         e.Degree,
		"field_of_study": e.FieldOfStudy,
		"description":    e.Description,
		"start_date":     e.StartDate,
		"end_date":       e.EndDate,
		"is_current":     e.IsCurrent,
		"updated_at":     pgutil.ToTimestampt(time.Now()),
	}
	q, args, err := psql.Update("educations").
		SetMap(setMap).
		Where(sq.Eq{"id": e.ID}).
		ToSql()
	if err != nil {
		return err
	}
	result, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("education with that id not found")
	}
	return nil
}

func (r *EducationRepoImpl) Delete(ctx context.Context, id domain.EducationID) error {
	q, args, err := psql.Delete("educations").
		Where(sq.Eq{"id": pgutil.ToUUID(uuid.UUID(id))}).
		ToSql()
	if err != nil {
		return err
	}
	result, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("education with that id not found")
	}
	return nil
}
