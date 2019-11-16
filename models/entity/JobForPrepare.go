package entity

import (
	"github.com/PuerkitoBio/goquery"
)

type JobForPrepare struct {
	Job goquery.Selection
	Location
	Page int
}
