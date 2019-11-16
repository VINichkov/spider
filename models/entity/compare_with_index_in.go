package entity

type CompareWithIndexIn struct {
	CompanyName string
	TitleJob string
	Url string
	Location
}

func NewCompareWithIndexIn(campany string, title string, url string, loc Location) *CompareWithIndexIn{
	return &CompareWithIndexIn{
		campany,
		title,
		url,
		loc,
	}
}
