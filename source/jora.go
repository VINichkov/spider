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




type Jora struct {
	Host string
	Url string
	ST string
	SP string

}


func InitJora() *Jora  {
	return &Jora{
		"https://au.jora.com",
		"https://au.jora.com/j?",
		"date",
		"facet_listed_date",
	}
}

func (source *Jora)Name() string{
	return "Jora"
}

func (source *Jora)CreateQuery(location entity.Location, page int) string{
	a :="24h"
	button := ""
	return fmt.Sprintf(
		"%sa=%s&button=%s&l=%s&p=%d&sp=%s&st=%s",
		source.Url, a, button, 	location.Suburb, page, 	source.SP, source.ST,
		)
}

func (source *Jora)NumberOfAds(doc goquery.Document) (int, error) {
	text :=doc.Find("div#search_info span,  div.search-count span").Last().Text()
	result := strings.Replace(text, ",","", -1)
	i , err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (source *Jora)GetTitle(job goquery.Selection) Title {
	tag := job.Find("a.jobtitle, a.job").First()
	href, _ := tag.Attr("href")
		return Title{
		href,
		tag.Text(),
	}
}

func (source *Jora)AgeOfJob(job goquery.Selection) string  {
	return job.Find("span.date").Text()
}

func (source *Jora)UrlToJob(url string) (string, error){
	idx := strings.IndexByte(url, '?');
	if  idx >= 0 {
		return  source.Host + url[:idx], nil
	}
	return "", errors.New("URL не найден")
}

func (source *Jora)Company(job goquery.Selection) string{
	return job.Find("div span.company").First().Text()
}

func (source *Jora)Salary(job goquery.Selection) []int {
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

func (source *Jora)ApplyLinc(job goquery.Document, obj entity.GotJob) string{
	//Ищим кнопку apply и получаем ссылка
	apply_linc, flag := job.Find("a.apply_link").First().Attr("href")

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
		log.Error().Msg(err.Error())
		return  obj.Link
	}

	return js.Url
}

func (source *Jora)Description(job goquery.Document) goquery.Selection {
	return *job.Find("div.summary").First()

}

