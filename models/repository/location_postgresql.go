package repository

import (
	"github.com/jmoiron/sqlx"
	"spider/models"
	"spider/models/entity"
)

func NewSQLLocationRepo(Conn *sqlx.DB) models.LocationRepository{
	return &postgresLocationRepo{
		Conn: Conn,
	}
}

type postgresLocationRepo struct {
	Conn *sqlx.DB
}

func (l *postgresLocationRepo) FindByID(id int)(*entity.Location, error){
	query := "SELECT ID, STATE, SUBURB FROM locations WHERE id = $1;"
	location := entity.Location{}
	err := l.Conn.Get(&location, query, id)
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (l *postgresLocationRepo) All()(*[]entity.Location, error){
	query := "SELECT ID, STATE, SUBURB FROM locations;"
	locations := []entity.Location{}
	err := l.Conn.Select(&locations, query,)
	if err != nil {
		return nil, err
	}
	return &locations, nil
}


