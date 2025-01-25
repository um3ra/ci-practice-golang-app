package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/um3ra/auth-microservice/internal/model"
	"github.com/um3ra/auth-microservice/internal/repository/converter"
	repoMod "github.com/um3ra/auth-microservice/internal/repository/model"
	"github.com/georgysavva/scany/v2/pgxscan"
)

const (
	tableName       = "users"
	limit           = 10
	nameColumn      = "name"
	emailColumn     = "email"
	idColumn        = "id"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

var selectQ = sq.Select(idColumn, nameColumn, emailColumn, updatedAtColumn, createdAtColumn).From(tableName)

type userRepository struct {
	pgPool *pgxpool.Pool
}

func NewUserRepository(pgPool *pgxpool.Pool) *userRepository {
	return &userRepository{
		pgPool: pgPool,
	}
}

func (r *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, updatedAtColumn, createdAtColumn).From(tableName).Limit(limit)
	sql, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	q, err := r.pgPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	var user repoMod.User

	users := make([]model.User, 0, limit)
	for q.Next() {
		err := q.Scan(&user.Id, &user.Name, &user.Email, &user.UpdatedAt, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		model := converter.ToUserFromRepo(&user)
		users = append(users, *model)
	}
	return users, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) (int64, error) {
	return 0, nil
}

func (r *userRepository) GetById(ctx context.Context, id int64) (*model.User, error) {
	builder := selectQ.PlaceholderFormat(sq.Dollar).Where("id=$1", id)
	sql, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}
	var user repoMod.User
	rows, err := r.pgPool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanOne(&user, rows); err != nil {
		return nil, err
	}
	return converter.ToUserFromRepo(&user), nil
}
