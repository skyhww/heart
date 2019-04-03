package entity

import (
	"time"
	"github.com/jmoiron/sqlx"
	"database/sql"
)

type Store struct {
	Id         int64      `db:"id" json:"id"`
	Url        *string    `db:"url" json:"url"`
	StoreType  *string    `db:"store_type" json:"-"`
	CreateTime *time.Time `db:"create_time" json:"create_time"`
	Suffix     string     `db:"suffix" json:"suffix"`
	Enable     int        `db:"enable" json:"-"`
}

type StorePersist interface {
	Save(store *Store) error
	Get(url string) (*Store, error)
}
type StoreDao struct {
	db *sqlx.DB
}

func NewStorePersist(db *sqlx.DB) StorePersist {
	return &StoreDao{db: db}
}

func (storeDao *StoreDao) Save(store *Store) error {
	tx := storeDao.db.MustBegin()
	r, err := tx.NamedExec("INSERT INTO store (url, store_type,suffix,enable,create_time) VALUES (:url ,:store_type,:suffix,1,:create_time)", store)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	id, _ := r.LastInsertId()
	store.Id = id
	return nil
}
func (storeDao *StoreDao) Get(url string) (*Store, error) {
	user := &Store{}
	err := storeDao.db.Get(user, "select url, store_type,suffix,create_time from store where url=? ", url)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}
