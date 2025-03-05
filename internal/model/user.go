package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	UpdatedAt sql.NullTime
	CreatedAt time.Time
}
