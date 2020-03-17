package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"spider/models"
	"spider/models/entity/client"
	"time"
)

func NewSQLdbClientRepo(Conn *sqlx.DB) models.ClientRepository {
	return &dbClientRepo{
		Conn: Conn,
	}
}

type dbClientRepo struct {
	Conn *sqlx.DB
}

func (l *dbClientRepo)FindByCompanyIdFirst(company_id int) (int, error){
	query := "SELECT ID FROM clients  WHERE company_id = $1 limit 1"
	var id int
	err := l.Conn.Get(&id, query, company_id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (l *dbClientRepo)Create(client *client.SimulationClient)(int, error){
	stmt, err := l.Conn.PrepareNamed(`INSERT INTO "clients" ("firstname", "lastname", "email", "location_id", 
				"created_at", "updated_at", "confirmed_at", "send_email", "alert", "company_id", "character") 
				VALUES (:firstname, :lastname, :email, :location_id, :created_at, :updated_at, :confirmed_at, :send_email,
				 :alert, :company_id, :character) RETURNING "id"`)
	if err != nil {
		return 0, err
	}

	timeNow := time.Now()
	client.CreatedAt = sql.NullTime{timeNow, true}
	client.UpdatedAt = sql.NullTime{timeNow, true}
	client.ConfirmedAt = sql.NullTime{timeNow, true}
	var id int
	err = stmt.Get(&id, client)
	if err != nil {
		return 0, err
	}

	return id, err

}