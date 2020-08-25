package repository

import (
	"github.com/jmoiron/sqlx"
	"spider/models"
)

func NewSQLdbIndustryRepo(Conn *sqlx.DB) models.IndustryRepository {
	return &dbIndustryRepo{
		Conn: Conn,
	}
}

type dbIndustryRepo struct {
	Conn *sqlx.DB
}

func (l *dbIndustryRepo)Other() (int, error) {
	query := "SELECT ID FROM industries where name = 'Other'"
	var id int
	err := l.Conn.Get(&id, query)
	if err != nil {
		return 0, err
	}
	return id, nil
}