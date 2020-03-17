package client

import (
	"database/sql"
	"spider/models/enum"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

type SimulationClient struct {
	FirstName sql.NullString`db:"firstname"`
	LastName sql.NullString `db:"lastname"`
	Email sql.NullString `db:"email"`
	LocationId sql.NullInt32 `db:"location_id"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	ConfirmedAt sql.NullTime `db:"confirmed_at"`
	SendEmail sql.NullBool `db:"send_email"`
	Alert sql.NullBool `db:"alert"`
	CompanyId sql.NullInt32 `db:"company_id"`
	Character string `db:"character"`
}

func NewSimulationClientForClawler(company_name string, company_id int, location_id int) *SimulationClient{
	const last_name = "HR"
	return &SimulationClient{
		FirstName:   sql.NullString{company_name, true},
		LastName:    sql.NullString{last_name, true},
		Email:       sql.NullString{emailGenerate(company_name), true},
		LocationId:  sql.NullInt32{int32(location_id), true},
		CreatedAt:   sql.NullTime{time.Now(), false},
		UpdatedAt:   sql.NullTime{time.Now(), false},
		ConfirmedAt: sql.NullTime{time.Now(), false},
		SendEmail:   sql.NullBool{false, true},
		Alert:   sql.NullBool{false, true},
		CompanyId:   sql.NullInt32{int32(company_id), true},
		Character:   enum.Employer(),
	}
}

func emailGenerate(company string) string{
	const domain ="@mail.com.au"
	//Массив вимволов для случайного выбора
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	//Размер строки
	const length = 5
	//Создаем массив нужного размера
	b := make([]byte, length)

	//Получаем случайное число сдвига
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	//Вставляем случайный символ, зависящий от времени
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	//Из названия компании удаляем лишник символы
	reg1 := regexp.MustCompile("\\s")
	reg2 := regexp.MustCompile("\\W")
	replaceStr := reg1.ReplaceAllString(company, "_")
	replaceStr = reg2.ReplaceAllString(replaceStr, "")
	//Собираем email
	return replaceStr + string(b) + domain
}


func (client *SimulationClient)String()string{
	str := ""
	if client.FirstName.Valid {str += "FirstName: " + client.FirstName.String + "\n" } else {str += "FirstName: null \n"}
	if client.LastName.Valid {str += "LastName: " + client.LastName.String + "\n" } else {str += "LastName: null \n"}
	if client.Email.Valid {str += "Email: " + client.Email.String + "\n" } else {str += "Email: null \n"}
	if client.LocationId.Valid {str += "LocationId: " +  strconv.Itoa(int(client.LocationId.Int32)) + "\n" } else {str +="LocationId: null\n"}
	if client.CreatedAt.Valid {str += "CreatedAt: " +  client.CreatedAt.Time.String() + "\n"} else {str +="CreatedAt: null\n"}
	if client.UpdatedAt.Valid {str += "UpdatedAt: " +  client.UpdatedAt.Time.String() + "\n"} else {str +="UpdatedAt: null\n"}
	if client.ConfirmedAt.Valid {str += "ConfirmedAt: " +  client.ConfirmedAt.Time.String() + "\n"} else {str +="ConfirmedAt: null\n"}
	if client.SendEmail.Valid {str += "SendEmail: " +  strconv.FormatBool(client.SendEmail.Bool) + "\n"} else {str +="SendEmail: null\n"}
	if client.Alert.Valid {str += "Alert: " +  strconv.FormatBool(client.Alert.Bool) + "\n"} else {str +="Alert: null\n"}
	if client.CompanyId.Valid {str += "CompanyId: " +  strconv.Itoa(int(client.CompanyId.Int32)) + "\n"} else {str +="CompanyId: null\n"}
	return str
}



