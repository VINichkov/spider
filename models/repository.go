package models

import (
	"spider/models/entity"
	"spider/models/entity/client"
	"spider/models/entity/company"
	"spider/models/entity/job"
)


type LocationRepository interface {
	FindByID(int)(*entity.Location, error)
	All()(*[]entity.Location, error)
}

type ClientRepository interface {
	FindByCompanyIdFirst(int) (int, error)
	Create( *client.SimulationClient)(int, error)
}

type SizeRepository interface {
	First() (int, error)
}

type IndustryRepository interface {
	Other() (int, error)
}

type CompanyRepository interface {
	FindCompanyByName(string)(*entity.CompanyId, error)
	Create( *company.Company)(int, error)
}

type JobRepository interface {
	FindJobBySourceOrTitle(int, int, string, string)(*entity.JobForCrawler, error)
	FindById(int)(*job.Job, error)
	Create( *job.Job)(int, error)
}