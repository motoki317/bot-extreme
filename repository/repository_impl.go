package repository

import "github.com/jmoiron/sqlx"

// Repository実装
type RepositoryImpl struct {
	db *sqlx.DB
}

func NewRepositoryImpl(db *sqlx.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}
