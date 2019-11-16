package source

import (
	"github.com/PuerkitoBio/goquery"
	"spider/models/entity"
)



type Source interface {
	Name() string
	CreateQuery(entity.Location, int) string
    NumberOfAds(goquery.Document) (int, error)
	GetTitle(goquery.Selection) Title
	AgeOfJob(goquery.Selection) string
	UrlToJob(string) (string, error)
	Company(goquery.Selection) string
	Salary(goquery.Selection) []int
	ApplyLinc(goquery.Document, entity.GotJob) string
	Description( goquery.Document) goquery.Selection
}
