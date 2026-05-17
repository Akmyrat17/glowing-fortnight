package persistence

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourorg/boilerplate/internal/domain"
	"github.com/yourorg/boilerplate/internal/modules/user/infra/persistence/dao"
	"github.com/yourorg/boilerplate/internal/shared/app_errors"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type UserRepoImpl struct {
	db *pgxpool.Pool
}

func NewUserRepoImpl(db *pgxpool.Pool) *UserRepoImpl {
	return &UserRepoImpl{db: db}
}

func (r *UserRepoImpl) Save(ctx context.Context, user *domain.User) error {
	d := dao.FromDomain(user)
	query, args, err := psql.Insert("users").
		Columns("id", "name", "email", "phone", "role", "password_hash", "token_version", "status", "created_at", "updated_at").
		Values(d.ID, d.Name, d.Email, d.Phone, d.Role, d.PasswordHash, d.TokenVersion, d.Status, d.CreatedAt, d.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return app_errors.DatabaseFailure(err)
	}
	return nil
}

func (r *UserRepoImpl) FindByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	query, args, err := psql.Select(
		"id", "name", "email", "phone", "role", "password_hash", "token_version", "status", "created_at", "updated_at",
	).From("users").Where(sq.Eq{"id": uuid.UUID(id)}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var d dao.UserDAO
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&d.ID, &d.Name, &d.Email, &d.Phone, &d.Role, &d.PasswordHash, &d.TokenVersion, &d.Status, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, app_errors.NotFound("user")
		}
		return nil, app_errors.DatabaseFailure(err)
	}

	return d.ToDomain(), nil
}

func (r *UserRepoImpl) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, int64, error) {
	query, args, err := psql.Select(
		"id", "name", "email", "phone", "role", "password_hash", "token_version", "status", "created_at", "updated_at",
	).From("users").Limit(uint64(limit)).Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, app_errors.DatabaseFailure(err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var d dao.UserDAO
		err = rows.Scan(&d.ID, &d.Name, &d.Email, &d.Phone, &d.Role, &d.PasswordHash, &d.TokenVersion, &d.Status, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, 0, app_errors.DatabaseFailure(err)
		}
		users = append(users, d.ToDomain())
	}

	if err = rows.Err(); err != nil {
		return nil, 0, app_errors.DatabaseFailure(err)
	}

	// Get total count
	countQuery, countArgs, err := psql.Select("COUNT(*)").From("users").ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var total int64
	err = r.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, app_errors.DatabaseFailure(err)
	}

	return users, total, nil
}

func (r *UserRepoImpl) Update(ctx context.Context, user *domain.User) error {
	d := dao.FromDomain(user)
	query, args, err := psql.Update("users").
		Set("name", d.Name).
		Set("email", d.Email).
		Set("phone", d.Phone).
		Set("role", d.Role).
		Set("status", d.Status).
		Set("token_version", d.TokenVersion).
		Set("updated_at", d.UpdatedAt).
		Where(sq.Eq{"id": d.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return app_errors.DatabaseFailure(err)
	}
	return nil
}

func (r *UserRepoImpl) Delete(ctx context.Context, id domain.UserID) error {
	query, args, err := psql.Delete("users").Where(sq.Eq{"id": uuid.UUID(id)}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return app_errors.DatabaseFailure(err)
	}
	return nil
}

func (r *UserRepoImpl) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query, args, err := psql.Select(
		"id", "name", "email", "phone", "role", "password_hash", "token_version", "status", "created_at", "updated_at",
	).From("users").Where(sq.Eq{"email": email}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var d dao.UserDAO
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&d.ID, &d.Name, &d.Email, &d.Phone, &d.Role, &d.PasswordHash, &d.TokenVersion, &d.Status, &d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, app_errors.NotFound("user")
		}
		return nil, app_errors.DatabaseFailure(err)
	}

	return d.ToDomain(), nil
}
