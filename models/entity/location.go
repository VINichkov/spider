package entity

type Location struct {
	Id int `db:"id"`
	State string `db:"state"`
	Suburb string `db:"suburb"`
	//нет части типов. Расширить по мере надобности
}
