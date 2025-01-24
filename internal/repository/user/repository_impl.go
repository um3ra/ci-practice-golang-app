package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/um3ra/auth-microservice/internal/repository/converter"
	repoMod "github.com/um3ra/auth-microservice/internal/repository/model"
	"github.com/um3ra/auth-microservice/internal/model"
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
	var users []model.User
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

func (r *userRepository) Create(ctx context.Context, user *model.User) (int64, error){
	return 0, nil
}

func (r *userRepository) GetById(ctx context.Context, id int64) (*model.User, error){
	return nil, nil
}