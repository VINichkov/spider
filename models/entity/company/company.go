package company

import (
	"database/sql"
	"strconv"
	"time"
)

type Company struct {
	Id sql.NullInt32 `db:"id"`
	Name sql.NullString`db:"name"`
	SizeId sql.NullInt32 `db:"size_id"`
	LocationId sql.NullInt32 `db:"location_id"`
	Site sql.NullString `db:"site"`
	LogoUID sql.NullString `db:"logo_uid"`
	RecruitmentAgency sql.NullBool `db:"recrutmentagency"`
	Description sql.NullString `db:"description"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	IndustryId sql.NullInt32 `db:"industry_id"`
	Big sql.NullBool `db:"big"`
	Names sql.NullString`db:"names"`
}

func NewCompany(company_name string, size_id int, location_id int, site string, logo_uid string,
	recruitment_agency bool, description string, industry_id int, big bool, names string) *Company{
	const last_name = "HR"
	return &Company{
		Id:   sql.NullInt32{0, false},
		Name: sql.NullString{company_name, company_name != ""},
		SizeId: sql.NullInt32{int32(size_id), size_id != 0},
		LocationId: sql.NullInt32{int32(location_id), size_id != 0},
		Site: sql.NullString{site, site != ""},
		LogoUID: sql.NullString{logo_uid, logo_uid != ""},
		RecruitmentAgency: sql.NullBool{recruitment_agency, true},
		Description: sql.NullString{description, description != ""},
		CreatedAt: sql.NullTime{time.Now(), false},
		UpdatedAt: sql.NullTime{time.Now(), false},
		IndustryId: sql.NullInt32{int32(industry_id), industry_id != 0},
		Big: sql.NullBool{big, true},
		Names: sql.NullString{names, names != ""},
	}
}

func (client *Company)String()string{
	str := ""
	if client.Id.Valid {str += "Id: " +  strconv.Itoa(int(client.Id.Int32)) + "\n" } else {str +="Id: null\n"}
	if client.Name.Valid {str += "Name: " + client.Name.String + "\n" } else {str += "Name: null \n"}
	if client.SizeId.Valid {str += "SizeId: " +  strconv.Itoa(int(client.SizeId.Int32)) + "\n" } else {str +="SizeId: null\n"}
	if client.LocationId.Valid {str += "LocationId: " +  strconv.Itoa(int(client.LocationId.Int32)) + "\n" } else {str +="LocationId: null\n"}
	if client.Site.Valid {str += "Site: " + client.Site.String + "\n" } else {str += "Site: null \n"}
	if client.LogoUID.Valid {str += "LogoUID: " + client.LogoUID.String + "\n" } else {str += "LogoUID: null \n"}
	if client.RecruitmentAgency.Valid {str += "RecruitmentAgency: " +  strconv.FormatBool(client.RecruitmentAgency.Bool) + "\n"} else {str +="SendEmail: null\n"}
	if client.Description.Valid {str += "Description: " + client.Description.String + "\n" } else {str += "Description: null \n"}
	if client.CreatedAt.Valid {str += "CreatedAt: " +  client.CreatedAt.Time.String() + "\n"} else {str +="CreatedAt: null\n"}
	if client.UpdatedAt.Valid {str += "UpdatedAt: " +  client.UpdatedAt.Time.String() + "\n"} else {str +="UpdatedAt: null\n"}
	if client.Big.Valid {str += "Big: " +  strconv.FormatBool(client.Big.Bool) + "\n"} else {str +="Big: null\n"}
	if client.Names.Valid {str += "Names: " + client.Names.String + "\n" } else {str += "Names: null \n"}
	return str
}