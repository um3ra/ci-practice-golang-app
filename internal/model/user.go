package model

import (
	"database/sql"
	"time"
)
type User struct {
	Id int64
	Name string
	Email string
	Password string
	UpdatedAt sql.NullTime
	CreatedAt time.Time
}