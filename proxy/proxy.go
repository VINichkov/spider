package proxy

import (
	"github.com/rs/zerolog/log"
	"net/http"
	"net/url"

	//"proxy/client"
	//"time"
)

type UrlFromRedirect struct {
	Body string `json:"body"`
	Url string 	`json:"uri"`
}

func Connect(arg string, mode bool) (*http.Response, error) {
	if arg == "" {
		arg = "https://rbc.ru"
	}

	var uri string

	// Если true то режим прокси
	// Иначе перулаем урл на новый сайт

	if mode {
		uri = "https://frozen-beach-74886.herokuapp.com/open?url=" + url.QueryEscape(arg)
	} else  {
		uri = "https://frozen-beach-74886.herokuapp.com/redirect?url=" + url.QueryEscape(arg)
	}
	//Локальный прокси
	/*if mode {
		uri = "http://localhost:8081/open?url=" + url.QueryEscape(arg)
	} else  {
		uri = "http://localhost:8081/redirect?url=" +  url.QueryEscape(arg)
	}*/

	//log.Debug().Msg(uri)

	return open(uri, 3)

}

func open(url string, index int8) (*http.Response, error)  {
	/*tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{
		Transport: tr,
		Timeout: 5 * time.Second,
	}*/
	resp, err := http.Get(url)
	log.Debug().Msg(url)
	/*if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		bodyString := string(bodyBytes)
		log.Info().Msg(bodyString)
	}*/
	if err != nil {
		if index > 0{
			return open(url, index -1)
		} else {
			return nil, err
		}
	}

	return resp, nil
}

