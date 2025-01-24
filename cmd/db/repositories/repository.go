package repositories

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	fake "github.com/brianvoe/gofakeit/v7"
	pg "github.com/jackc/pgx/v5"
)

type Repository struct {
	Db      *pg.Conn
	Context context.Context
}

type RepositoryDeps struct {
	Db      *pg.Conn
	Context context.Context
}

func NewRepository(deps RepositoryDeps) *Repository {
	return &Repository{
		Db:      deps.Db,
		Context: deps.Context,
	}
}

func (repo *Repository) Save() (bool, error) {
	defer repo.Db.Close(repo.Context)

	sql, args, err := sq.Insert("users").PlaceholderFormat(sq.Dollar).Columns("name", "email", "password").Values(fake.Name(), fake.Email(), fake.Name()).Suffix("RETURNING id").ToSql()
	// sql, _, err := sq.Insert("users").PlaceholderFormat(sq.Dollar).Columns("name", "email", "password").Values(fake.Name(), fake.Email(), fake.StreetName()).ToSql()
	res, err := repo.Db.Exec(repo.Context, sql, args...)
	return res.RowsAffected() > 0, err
}


type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}


func (repo *Repository) GetAll() ([]User, error){
	defer repo.Db.Close(repo.Context)

	sql, _, err := sq.Select("id", "name", "email").From("users").ToSql()
	if err != nil {
		return nil, err
	}

	var users []User

	res,err := repo.Db.Query(repo.Context, sql)
	var name string
	var id int
	var email string


	for res.Next(){
		err := res.Scan(&id, &name, &email);
		if err != nil{
			fmt.Println(err.Error())
			return nil, err
		}
		users = append(users, User{
			Id: id,
			Name: name,
			Email: email,
		})
	}

	return users, nil
}

// func (repo *Repository) GetById(id int) (*User, error) {

// 	defer repo.Db.Close(repo.Context)
// 	fmt.Println(id)

// 	sql, args, err := sq.Select("id", "name", "email").PlaceholderFormat(sq.Dollar).From("users").Where("id=$1", id).ToSql()
// 	fmt.Println(sql)

// 	if err != nil {
// 		return nil, err
// 	}

// 	res, err := repo.Db.Query(repo.Context, sql, args...)

// 	if err != nil {
// 		return nil, err
// 	}

// 	user := &User{}

// 	if err := res.Scan(&user.Id, &user.Name, &user.Email); err != nil {
// 		return nil, err
// 	}
// 	fmt.Printf("%#v\n", user)

// 	return user, nil
// }

func (repo *Repository) GetById(id int) (*User, error) {

	sql, args, err := sq.Select("id", "name", "email").
			PlaceholderFormat(sq.Dollar).
			From("users").
			Where("id=$1", id).
			ToSql()


	if err != nil {
			return nil, err
	}

	row := repo.Db.QueryRow(repo.Context, sql, args...)

	user := &User{}
	
	if err := row.Scan(&user.Id, &user.Name, &user.Email); err != nil {
			if err != nil {
					return nil, fmt.Errorf("user with id %d not found", id)
			}
			return nil, err
	}


	return user, nil
}
