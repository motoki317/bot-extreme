package repository

import (
	"github.com/jmoiron/sqlx"
	"sync"
)

// Repository実装
type RepositoryImpl struct {
	db          *sqlx.DB
	lock        sync.Mutex
	channelLock sync.Mutex
}

func NewRepositoryImpl(db *sqlx.DB) Repository {
	return &RepositoryImpl{
		db: db,
	}
}

func (r *RepositoryImpl) Lock() {
	r.lock.Lock()
}

func (r *RepositoryImpl) Unlock() {
	r.lock.Unlock()
}

func (r *RepositoryImpl) ChannelLock() {
	r.channelLock.Lock()
}

func (r *RepositoryImpl) ChannelUnlock() {
	r.channelLock.Unlock()
}
