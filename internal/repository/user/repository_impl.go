package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/um3ra/auth-microservice/internal/client/db"
	"github.com/um3ra/auth-microservice/internal/model"
	"github.com/um3ra/auth-microservice/internal/repository/user/converter"
	repoMod "github.com/um3ra/auth-microservice/internal/repository/user/model"
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
	dbClient db.Client
}

func NewUserRepository(client db.Client) *userRepository {
	return &userRepository{
		dbClient: client,
	}
}

func (r *userRepository) GetAll(ctx context.Context) ([]model.User, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn, updatedAtColumn, createdAtColumn).From(tableName).Limit(limit)
	sql, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}
	q, err := r.dbClient.DB().QueryContext(ctx, db.Query{QueryRaw: sql}, args...)
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
	builder := sq.Insert(tableName).PlaceholderFormat(sq.Dollar).Columns(nameColumn, emailColumn, "password").Values(user.Name, user.Email, user.Password).Suffix("RETURNING id")
	sql, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}
	var id int64

	err = r.dbClient.DB().QueryRowContext(ctx, db.Query{QueryRaw: sql, Name: "user creating sql req"}, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *userRepository) GetById(ctx context.Context, id int64) (*model.User, error) {
	builder := selectQ.PlaceholderFormat(sq.Dollar).Where("id=$1", id)
	sql, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}
	var user repoMod.User
	rows, err := r.dbClient.DB().QueryContext(ctx, db.Query{QueryRaw: sql}, args...)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanOne(&user, rows); err != nil {
		return nil, err
	}
	return converter.ToUserFromRepo(&user), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	builder := selectQ.PlaceholderFormat(sq.Dollar).Where("email=$1", email)
	sql, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}
	var user repoMod.User
	rows, err := r.dbClient.DB().QueryContext(ctx, db.Query{QueryRaw: sql}, args...)
	if err != nil {
		return nil, err
	}
	if err := pgxscan.ScanOne(&user, rows); err != nil {
		return nil, err
	}
	return converter.ToUserFromRepo(&user), nil
}
