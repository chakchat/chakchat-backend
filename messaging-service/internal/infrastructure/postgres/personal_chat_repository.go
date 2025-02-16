package postgres

import (
	"database/sql"
)

type PersonalChatRepository struct {
	db *sql.DB
}

func NewPersonalChatRepository(db *sql.DB) *PersonalChatRepository {
	return &PersonalChatRepository{
		db: db,
	}
}
