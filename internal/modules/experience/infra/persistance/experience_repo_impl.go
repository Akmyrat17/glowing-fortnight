package persistance

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/experience/infra/persistance/dao"
	projectDAO "github.com/boilerplate/internal/modules/projects/infra/persistance/dao"
	"github.com/boilerplate/internal/shared/pgutil"
	"github.com/boilerplate/pkg/query"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var experienceAllowedFields = map[string]string{
	"id":         "id",
	"company":    "company",
	"position":   "position",
	"is_current": "is_current",
	"location":   "location",
	"start_date": "start_date",
	"end_date":   "end_date",
	"created_at": "created_at",
	"updated_at": "updated_at",
}

var experienceColumns = []string{
	"e.id", "e.company", "e.position", "e.description",
	"e.start_date", "e.end_date", "e.is_current", "e.location",
	"e.created_at", "e.updated_at", "e.project_id",
	// project fields — null if no project linked
	"p.id", "p.name", "p.url", "p.repo_url",
	"p.status", "p.project_type", "p.tags",
}

var experienceUpsertColumns = []string{
	"id", "company", "position", "description",
	"start_date", "end_date", "is_current", "location",
	"created_at", "updated_at",
}

func scanExperienceWithProject(row pgx.Row) (*dao.ExperienceDAO, *projectDAO.ProjectDAO, error) {
	var e dao.ExperienceDAO
	var p projectDAO.ProjectDAO
	// p fields are nullable since it's a LEFT JOIN
	var pID, pName, pUrl, pRepoUrl, pStatus, pProjectType pgtype.Text
	var pTags pgtype.Array[string]

	err := row.Scan(
		&e.ID, &e.Company, &e.Position, &e.Description,
		&e.StartDate, &e.EndDate, &e.IsCurrent, &e.Location,
		&e.CreatedAt, &e.UpdatedAt, &e.ProjectID,
		&pID, &pName, &pUrl, &pRepoUrl, &pStatus, &pProjectType, &pTags,
	)
	if err != nil {
		return nil, nil, err
	}

	if pID.Valid {
		p.ID = pgutil.ToUUID(uuid.MustParse(pID.String))
		p.Name = pName.String
		p.Url = pUrl
		p.RepoUrl = pRepoUrl
		p.Status = pStatus.String
		p.ProjectType = pProjectType.String
		p.Tags = pTags
		return &e, &p, nil
	}

	return &e, nil, nil
}

type ExperienceRepoImpl struct {
	db *pgxpool.Pool
}

func NewExperienceRepoImpl(db *pgxpool.Pool) *ExperienceRepoImpl {
	return &ExperienceRepoImpl{db: db}
}

func (r *ExperienceRepoImpl) Create(ctx context.Context, e *dao.ExperienceDAO) error {
	cols := append([]string{}, experienceUpsertColumns...)
	vals := []interface{}{
		e.ID, e.Company, e.Position, e.Description,
		e.StartDate, e.EndDate, e.IsCurrent, e.Location,
		e.CreatedAt, e.UpdatedAt,
	}

	if e.ProjectID.Valid {
		cols = append(cols, "project_id")
		vals = append(vals, e.ProjectID)
	}

	q, args, err := psql.Insert("experiences").
		Columns(cols...).
		Values(vals...).
		ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, q, args...)
	return err
}

func (r *ExperienceRepoImpl) FindByID(ctx context.Context, id domain.ExperienceID) (*domain.Experience, error) {
	q, args, err := psql.Select(experienceColumns...).
		From("experiences e").
		LeftJoin("projects p ON p.id = e.project_id").
		Where(sq.Eq{"e.id": pgutil.ToUUID(uuid.UUID(id))}).
		ToSql()
	if err != nil {
		return nil, err
	}

	e, p, err := scanExperienceWithProject(r.db.QueryRow(ctx, q, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("experience with that id not found")
		}
		return nil, err
	}

	exp := e.ToDomain()
	if p != nil {
		exp.Project = p.ToDomain()
	}
	return exp, nil
}

func (r *ExperienceRepoImpl) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Experience, int64, error) {
	countBuilder := psql.Select("COUNT(*)").From("experiences e").LeftJoin("projects p ON p.id = e.project_id")
	countBuilder, err := query.ApplyFilters(countBuilder, filters, experienceAllowedFields)
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

	qb := psql.Select(experienceColumns...).From("experiences e").LeftJoin("projects p ON p.id = e.project_id")
	qb, err = query.ApplyFilters(qb, filters, experienceAllowedFields)
	if err != nil {
		return nil, 0, err
	}
	if len(sorts) > 0 {
		qb = query.ApplySort(qb, sorts, experienceAllowedFields)
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

	var result []*domain.Experience
	for rows.Next() {
		e, p, err := scanExperienceWithProject(rows)
		if err != nil {
			return nil, 0, err
		}
		exp := e.ToDomain()
		if p != nil {
			exp.Project = p.ToDomain()
		}
		result = append(result, exp)
	}
	return result, total, nil
}

func (r *ExperienceRepoImpl) Update(ctx context.Context, e *dao.ExperienceDAO) error {
	setMap := sq.Eq{
		"company":     sq.Expr("COALESCE(NULLIF(?::text, ''), company)", e.Company),
		"position":    sq.Expr("COALESCE(NULLIF(?::text, ''), position)", e.Position),
		"description": e.Description,
		"start_date":  e.StartDate,
		"end_date":    e.EndDate,
		"is_current":  e.IsCurrent,
		"location":    e.Location,
		"updated_at":  pgutil.ToTimestampt(time.Now()),
	}

	if e.ProjectID.Valid {
		setMap["project_id"] = e.ProjectID
	}

	q, args, err := psql.Update("experiences").
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
		return fmt.Errorf("experience with that id not found")
	}
	return nil
}

func (r *ExperienceRepoImpl) Delete(ctx context.Context, id domain.ExperienceID) error {
	q, args, err := psql.Delete("experiences").
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
		return fmt.Errorf("experience with that id not found")
	}
	return nil
}
