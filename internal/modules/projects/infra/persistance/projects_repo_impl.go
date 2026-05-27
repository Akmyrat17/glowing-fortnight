package persistance

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/projects/infra/persistance/dao"
	"github.com/boilerplate/internal/shared/pgutil"
	"github.com/boilerplate/pkg/query"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

var projectAllowedFields = map[string]string{
	"id":           "id",
	"name":         "name",
	"description":  "description",
	"url":          "url",
	"repo_url":     "repo_url",
	"start_date":   "start_date",
	"end_date":     "end_date",
	"status":       "status",
	"project_type": "project_type",
	"tags":         "tags",
	"created_at":   "created_at",
	"updated_at":   "updated_at",
}

var projectColumns = []string{
	"id", "name", "description", "url", "repo_url", "start_date", "end_date", "status", "project_type", "tags", "created_at", "updated_at",
}

func scanProject(row pgx.Row) (*dao.ProjectDAO, error) {
	var project dao.ProjectDAO
	err := row.Scan(
		&project.ID, &project.Name, &project.Description,
		&project.Url, &project.RepoUrl, &project.StartDate, &project.EndDate,
		&project.Status, &project.ProjectType, &project.Tags,
		&project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

type ProjectRepoImpl struct {
	db *pgxpool.Pool
}

func NewProjectRepoImpl(db *pgxpool.Pool) *ProjectRepoImpl {
	return &ProjectRepoImpl{db: db}
}

func (repo *ProjectRepoImpl) Create(ctx context.Context, dao *dao.ProjectDAO) error {
	query, args, err := psql.Insert("projects").
		Columns(projectColumns...).
		Values(dao.ID, dao.Name, dao.Description,
			dao.Url, dao.RepoUrl,
			dao.StartDate, dao.EndDate,
			dao.Status, dao.ProjectType,
			dao.Tags,
			dao.CreatedAt, dao.UpdatedAt,
		).
		ToSql()
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(ctx, query, args...)
	return err
}

func (repo *ProjectRepoImpl) FindByID(ctx context.Context, id domain.ProjectID) (*domain.Project, error) {
	query, args, err := psql.Select(projectColumns...).
		From("projects").
		Where(sq.Eq{"id": pgutil.ToUUID(uuid.UUID(id))}).
		ToSql()
	if err != nil {
		return nil, nil
	}
	projectDAO, err := scanProject(repo.db.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("project with that id not found")
		}
		return nil, err
	}
	return projectDAO.ToDomain(), nil
}

func (repo *ProjectRepoImpl) FindAll(ctx context.Context, limit, offset int, filters []query.Filter, sorts []query.SortField) ([]*domain.Project, int64, error) {
	countBuilder := psql.Select("COUNT(*)").From("projects")
	countBuilder, err := query.ApplyFilters(countBuilder, filters, projectAllowedFields)
	if err != nil {
		return nil, 0, err
	}
	countQuery, countArgs, err := countBuilder.ToSql()
	if err != nil {
		return nil, 0, err
	}
	var total int64
	if err := repo.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryBuilder := psql.Select(projectColumns...).From("projects")
	queryBuilder, err = query.ApplyFilters(queryBuilder, filters, projectAllowedFields)
	if err != nil {
		return nil, 0, err
	}
	if len(sorts) > 0 {
		queryBuilder = query.ApplySort(queryBuilder, sorts, projectAllowedFields)
	} else {
		queryBuilder = queryBuilder.OrderBy("created_at DESC")
	}
	dbQuery, args, err := queryBuilder.Limit(uint64(limit)).Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, 0, err
	}
	rows, err := repo.db.Query(ctx, dbQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []*domain.Project

	for rows.Next() {
		projectDAO, err := scanProject(rows)
		if err != nil {
			return nil, 0, err
		}
		projects = append(projects, projectDAO.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return projects, total, nil

}

func (repo *ProjectRepoImpl) Update(ctx context.Context, project *dao.ProjectDAO) error {
	setMap := sq.Eq{
		// text fields — only update if non-empty, otherwise keep existing value
		"name":         sq.Expr("COALESCE(NULLIF(?::text, ''), name)", project.Name),
		"description":  sq.Expr("COALESCE(NULLIF(?::text, ''), description)", project.Description),
		"url":          sq.Expr("COALESCE(NULLIF(?::text, ''), url)", project.Url),
		"repo_url":     sq.Expr("COALESCE(NULLIF(?::text, ''), repo_url)", project.RepoUrl),
		"status":       sq.Expr("COALESCE(NULLIF(?::text, ''), status)", project.Status),
		"project_type": sq.Expr("COALESCE(NULLIF(?::text, ''), project_type)", project.ProjectType),

		// nullable dates — set directly, nil means "clear the date"
		"start_date": project.StartDate,
		"end_date":   project.EndDate,

		// arrays — set directly
		"tags": project.Tags,

		// always update this
		"updated_at": pgutil.ToTimestampt(time.Now()),
	}

	query, args, err := psql.Update("projects").
		SetMap(setMap).
		Where(sq.Eq{"id": project.ID}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := repo.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("project with that id not found")
	}

	return nil
}

func (repo *ProjectRepoImpl) Delete(ctx context.Context, id domain.ProjectID) error {
	query, args, err := psql.Delete("projects").
		Where(sq.Eq{"id": pgutil.ToUUID(uuid.UUID(id))}).
		ToSql()
	if err != nil {
		return err
	}

	result, err := repo.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("project with that id not found")
	}

	return nil
}
