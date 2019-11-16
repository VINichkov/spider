package repository

import (
	"github.com/jmoiron/sqlx"
	"spider/models"
)

func NewSQLdbSizeRepo(Conn *sqlx.DB) models.SizeRepository {
	return &dbSizeRepo{
		Conn: Conn,
	}
}

type dbSizeRepo struct {
	Conn *sqlx.DB
}

func (l *dbSizeRepo)First()(int, error){
	query := "SELECT ID FROM sizes limit 1"
	var id int
	err := l.Conn.Get(&id, query)
	if err != nil {
		return 0, err
	}
	return id, nil
}
