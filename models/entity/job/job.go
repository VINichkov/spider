package job

import (
	"database/sql"
	"errors"
	"strconv"
	"time"
)

type Job struct {
	Id sql.NullInt32 `db:"id"`
	Title sql.NullString`db:"title"`
	LocationID sql.NullInt32 `db:"location_id"`
	SalaryMin sql.NullInt32 `db:"salarymin"`
	SalaryMax sql.NullInt32 `db:"salarymax"`
	Description sql.NullString `db:"description"`
	CompanyId sql.NullInt32 `db:"company_id"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	HighLight sql.NullTime `db:"highlight"`
	Top sql.NullTime `db:"top"`
	Urgent sql.NullTime `db:"urgent"`
	ClientId sql.NullInt32 `db:"client_id"`
	Close  sql.NullTime `db:"close"`
	Twitter sql.NullString `db:"twitter"`
	Source sql.NullString `db:"sources"`
	Apply sql.NullString `db:"apply"`
	ViewedCount sql.NullInt32 `db:"viewed_count"`
}

func NewJob(title string, location_id int,  salaryMin int, salaryMax int, description string,
	company_id int, client_id int, twitter string, source string, apply string) *Job {
	//	var salaryMin , salary_Max sql.NullInt32


	return &Job{
			sql.NullInt32{0, false},
		sql.NullString{title, title != ""},
		sql.NullInt32{int32(location_id), location_id > 0},
		sql.NullInt32{int32(salaryMin), salaryMin > 0},
		sql.NullInt32{int32(salaryMax), salaryMax > 0},
		sql.NullString{description, description != ""},
		sql.NullInt32{int32(company_id), company_id > 0},
		sql.NullTime{time.Now(), false},
		sql.NullTime{time.Now(), false},
		sql.NullTime{time.Now(), false},
		sql.NullTime{time.Now(), false},
		sql.NullTime{time.Now(), false},
		sql.NullInt32{int32(client_id), client_id > 0},
		sql.NullTime{time.Now(), false},
		sql.NullString{twitter, twitter != ""},
		sql.NullString{source, source != ""},
		sql.NullString{apply, apply != ""},
		sql.NullInt32{0, true},
	}
}

func (job *Job)Valid()(bool, error){
	if job.Title.Valid == false{
		return false, errors.New("Title is empty")
	}

	if job.LocationID.Valid == false{
		return false, errors.New("Location is empty")
	}

	if job.CompanyId.Valid == false{
		return false, errors.New("Company is empty")
	}
	if job.ClientId.Valid == false{
		return false, errors.New("Client is empty")
	}

	return true, nil
}

func (job *Job)String()string{
	str := ""
	if job.Title.Valid {str += "title: " + job.Title.String + "\n" } else {str += "title: null \n"}
	if job.LocationID.Valid {str += "Location_id: " +  strconv.Itoa(int(job.LocationID.Int32)) + "\n" } else {str +="Location_id: null\n"}
	if job.SalaryMin.Valid {str += "SalaryMin: " +  strconv.Itoa(int(job.SalaryMin.Int32)) + "\n" } else {str +="SalaryMin: null\n"}
	if job.SalaryMax.Valid {str += "SalaryMax: " +  strconv.Itoa(int(job.SalaryMax.Int32)) + "\n" } else {str +="SalaryMax: null\n"}
	if job.Description.Valid {str += "Description: " +  job.Description.String + "\n" } else {str +="Description: null\n"}
	if job.CompanyId.Valid {str += "CompanyId: " +  strconv.Itoa(int(job.CompanyId.Int32)) + "\n"} else {str +="CompanyId: null\n"}
	if job.ClientId.Valid {str += "ClientId: " +  strconv.Itoa(int(job.ClientId.Int32)) + "\n"} else {str +="ClientId: null\n"}
	if job.CreatedAt.Valid {str += "CreatedAt: " +  job.CreatedAt.Time.String() + "\n"} else {str +="CreatedAt: null\n"}
	if job.UpdatedAt.Valid {str += "UpdatedAt: " +  job.UpdatedAt.Time.String() + "\n"} else {str +="UpdatedAt: null\n"}
	if job.HighLight.Valid {str += "HighLight: " +  job.HighLight.Time.String() + "\n"} else {str +="HighLight: null\n"}
	if job.Top.Valid {str += "Top: " +  job.Top.Time.String() + "\n"} else {str +="Top: null\n"}
	if job.Urgent.Valid {str += "Urgent: " +  job.Urgent.Time.String() + "\n"} else {str +="Urgent: null\n"}
	if job.Close.Valid {str += "Close: " +  job.Close.Time.String() + "\n"} else {str +="Close: null\n"}
	if job.Twitter.Valid {str += "Twitter: " + job.Twitter.String + "\n" } else {str += "Twitter: null \n"}
	if job.Source.Valid {str += "Source: " + job.Source.String + "\n" } else {str += "Source: null \n"}
	if job.Apply.Valid {str += "Apply: " + job.Apply.String + "\n" } else {str += "Apply: null \n"}
	if job.ViewedCount.Valid {str += "ViewedCount: " +  strconv.Itoa(int(job.ViewedCount.Int32)) + "\n"} else {str +="ViewedCount: null\n"}
	return str
}



