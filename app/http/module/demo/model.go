package demo

import (
	"database/sql"
	"time"
)

type UserModel struct {
	UserId int
	Name   string
	Age    int
}

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     *time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
