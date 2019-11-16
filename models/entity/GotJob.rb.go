package entity

import (
	"strconv"
)

type GotJob struct {
	Link string
	Title string
	Company string
	CompanyId
	Salary_min int
	Salary_max int
	Location
	Page int
}

func (job *GotJob)String() string{
	str := ""
	str += "title: " + job.Title + "\n"
	str += "Location_id: " +  strconv.Itoa(job.Location.Id) + "\n"
	str += "SalaryMin: " +  strconv.Itoa(job.Salary_min) + "\n"
	str += "SalaryMax: " +  strconv.Itoa(job.Salary_max) + "\n"
	str += "CompanyId: " +  strconv.Itoa(job.CompanyId.Id) + "\n"
	str += "Company: " + job.Company + "\n"
	str += "Source: " + job.Link + "\n"
	str += "Page: " +  strconv.Itoa(job.Page) + "\n"
	return str
}
