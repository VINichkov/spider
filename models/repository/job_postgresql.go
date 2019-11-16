package repository

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"spider/models"
	"spider/models/entity"
	"spider/models/entity/job"
	"strings"
	"time"
)

func NewSQLdbJobRepo(Conn *sqlx.DB) models.JobRepository {
	return &dbJobRepo{
		Conn: Conn,
	}
}

type dbJobRepo struct {
	Conn *sqlx.DB
}

func (l *dbJobRepo) FindJobBySourceOrTitle(location_id int, company_id int, source string, title string)(*entity.JobForCrawler, error){
	titleLower := strings.ToLower(title)
	sourceLower := strings.ToLower(source)
	query := "SELECT j.id, j.sources, j.title FROM jobs j WHERE j.location_id = $1 AND j.company_id = $2 AND (LOWER(j.sources) = $3 or LOWER(j.title) = $4);"
	job := entity.JobForCrawler{}
	err := l.Conn.Get(&job, query, location_id, company_id, sourceLower, titleLower)
	if models.EqualErrors(err, models.ErrNotFound){
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (l *dbJobRepo) FindById(id int)(*job.Job, error){
	query := `	SELECT 	id, 
						title, 
						location_id, 
						salarymin, 
						salarymax, 
						description,
						company_id,
						created_at,
						updated_at,
						highlight,
						top,
						urgent,
						client_id,
						close,
						twitter,
						sources,
						apply,
						viewed_count 
				FROM jobs 
				WHERE id = $1;`
	job := job.Job{}
	err := l.Conn.Get(&job, query, id)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (l *dbJobRepo) Create(job *job.Job)(int, error){
	stmt, err := l.Conn.PrepareNamed(`INSERT INTO "jobs" ("title", "location_id", "salarymin", "salarymax", 
				"description", "company_id", "created_at", "updated_at", "client_id", "close", "sources", "apply") 
				VALUES (:title, :location_id, :salarymin, :salarymax, :description, :company_id, :created_at, :updated_at,
				 :client_id, :close, :sources, :apply) RETURNING "id"`)
	if err != nil {
		return 0, err
	}

	timeNow := time.Now()
	job.CreatedAt = sql.NullTime{timeNow, true}
	job.UpdatedAt = sql.NullTime{timeNow, true}
	job.Close = sql.NullTime{timeNow.AddDate(0,0,15), true}
	var id int
	err = stmt.Get(&id, job)
	if err != nil {
		return 0, err
	}

	return id, err
}