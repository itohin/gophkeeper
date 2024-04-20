package entities

import (
	"database/sql"
)

const (
	TypeText = iota + 1
	TypePassword
)

type Secret struct {
	ID         string
	Name       string
	SecretType uint32
	Notes      string
	Data       interface{}
	UserID     string
	DeletedAt  sql.NullTime
}

type Text struct {
	Text string
}

type Password struct {
	Login    string
	Password string
}
