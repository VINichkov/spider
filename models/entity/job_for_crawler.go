package entity

type JobForCrawler struct {
	Id int `db:"id"`
	Source string `db:"sources"`
	Title string `db:"title"`
}
