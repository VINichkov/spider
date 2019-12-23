package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"spider/models/entity"
	"spider/proxy"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)


type Indeed struct {
	Host string
	Url string
	ST string

}


func InitIndeed() *Indeed  {
	return &Indeed{
		"https://au.indeed.com",
		"https://au.indeed.com/jobs?",
		"date",
	}
}
func (source *Indeed)Name() string{
	return "Indeed"
}

func (source *Indeed)CreateQuery(location entity.Location, page int) string{
	return fmt.Sprintf(
		"%sl=%s&sort=%s&start=%d&q=",
		source.Url,  location.Suburb, source.ST, page * 10,
	)
}

func (source *Indeed)NumberOfAds(doc goquery.Document) (int, error) {
	text :=doc.Find("div#searchCount").Last().Text()
	result := strings.Replace(text, ",","", -1)
	r, _ := regexp.Compile("\\d+")
	strs := r.FindAllString(result,-1)
	if len(strs)>0{
		i , err := strconv.Atoi(strs[len(strs)-1])
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, nil

}

func (source *Indeed)GetTitle(job goquery.Selection) Title {
	tag := job.Find("a.turnstileLink").First()
	href, _ := tag.Attr("href")
		return Title{
		href,
		tag.Text(),
	}
}

func (source *Indeed)AgeOfJob(job goquery.Selection) string  {
	return job.Find("div.result-link-bar span.date").Text()
}

func (source *Indeed)UrlToJob(url string) (string, error){
	r, _ := regexp.Compile("^/company")
	strs := r.FindAllString(url,-1)
	if len(strs)>0{
		return source.Host + url, nil
	}
	idx := strings.IndexByte(url, '?');
	if  idx > 0 {
		return  source.Host + "/viewjob" + url[idx:], nil
	}
	return "", errors.New("URL не найден")
}


func (source *Indeed)Company(job goquery.Selection) string{
	return job.Find("span.company").First().Text()
}

func (source *Indeed)Salary(job goquery.Selection) []int {
	str := job.Find("div.salary").Text()
	strAfteReplace := strings.Replace(str, ",", "", -1)
	r, _ := regexp.Compile("\\d+")
	strs := r.FindAllString(strAfteReplace,-1)
	result := []int{}
	for _, str := range strs {
		buf, _ := strconv.Atoi(str)
		result = append(result, buf)
	}
	return result
}

func (source *Indeed)ApplyLinc(job goquery.Document, obj entity.GotJob) string{
	//Ищим кнопку apply и получаем ссылка
	apply_linc, flag := job.Find("div#viewJobButtonLinkContainer a, div#jobsearch-ViewJobButtons-container a").First().Attr("href")

	//Если не нали кнопку или нет ссылки, то возвращаем линк на вакансию jora
	if !flag  {
		return obj.Link
	}

	//Проверяем URL
	url_from_button := apply_linc
	u, err := url.Parse(url_from_button)
	if err != nil {
		log.Error().Msg("Error :ApplyLinc: " + err.Error())
		return obj.Link
	}
	//Если нет хоста, то цеаляем его
	if u.Host == "" {
		url_from_button = source.Host + url_from_button
	}

	//Проверяем полученную ссылку
	new_url, err := proxy.Connect(url_from_button, false)

	if err != nil {
		log.Error().Msg("Error :ApplyLinc: Redirect is not work :" + err.Error())
		return obj.Link
	}

	defer new_url.Body.Close()

	js := &proxy.UrlFromRedirect{}
	jsonResp := json.NewDecoder(new_url.Body)
	err = jsonResp.Decode(js)

	if err != nil{
		log.Error().Msg("ApplyLinc: " + source.Name() + ": " +err.Error())
		return  obj.Link
	}

	return js.Url
}

func (source *Indeed)Description(job goquery.Document) goquery.Selection {
	return *job.Find("div#jobDescriptionText").First()
}

