package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"spider/models"
	"spider/models/entity"
	"spider/models/entity/company"
	"strings"
	"time"
)

func NewSQLCompanyRepo(Conn *sqlx.DB) models.CompanyRepository {
	return &dbCompanyRepo{
		Conn: Conn,
	}
}

type dbCompanyRepo struct {
	Conn *sqlx.DB
}

func (l *dbCompanyRepo) FindCompanyByName(name string)(*entity.CompanyId, error){
	nameLower := strings.ToLower(name)
	query := "SELECT id as \"result\" FROM companies WHERE LOWER(name) = $1 or ($1 = ANY(lower(names::text)::text[]));"
	companyId := entity.CompanyId{}
	err := l.Conn.Get(&companyId, query, nameLower)
	resEqualErrors := models.EqualErrors(err, models.ErrNotFound)

	if resEqualErrors{
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &companyId, nil
}

func (l *dbCompanyRepo)Create(company *company.Company)(int, error){
	stmt, err := l.Conn.PrepareNamed(`INSERT INTO "companies" ("name", "size_id", "location_id", "site", 
				"logo_uid", "recrutmentagency", "description", "created_at", "updated_at", "industry_id", "big", "names") 
				VALUES (:name, :size_id, :location_id, :site, :logo_uid, :recrutmentagency, :description, :created_at,
				 :updated_at, :industry_id, :big, :names) RETURNING "id"`)
	if err != nil {
		return 0, err
	}

	timeNow := time.Now()
	company.CreatedAt = sql.NullTime{timeNow, true}
	company.UpdatedAt = sql.NullTime{timeNow, true}

	if !company.IndustryId.Valid {
		repIndustry := NewSQLdbIndustryRepo(l.Conn)
		industryId, err := repIndustry.Other()

		if err != nil {
			return 0, err
		}

		company.IndustryId = sql.NullInt32 {int32(industryId), true}
	}

	var id int
	err = stmt.Get(&id, company)
	if err != nil {
		return 0, err
	}

	return id, err

}

