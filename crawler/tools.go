package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"regexp"
)

// Когда разместили объявление. Берем только сегоднешние
func how_long(text string) bool{
	var result bool
	res, err := regexp.MatchString(`hour|1 day|minute|Today`, text)
	if err != nil {
		res = false
	}
	if text == "Just posted" || text == "" || res {
		result =  true
	} else {
		result =  false
	}
	return result
}


func logging(function string, thread int, msg string){
	if function == ""{
		function = "NoName"
	}
	log.Error().Msg(fmt.Sprintf("%s %d: %s\n", function, thread, msg))
}

func loggingDB(function string, thread int, msg string){
	if function == ""{
		function = "NoName"
	}
	log.Info().Msg(fmt.Sprintf("%s %d: %s\n", function, thread, msg))
}

func deleteAttr(html *goquery.Selection){
	n := html.Nodes
	for _, elem := range n {
		for _, at :=range elem.Attr{
			html.RemoveAttr(at.Key)
		}
	}
}

func replace (str *string) string{
	rule1 := regexp.MustCompile("<h1>|<h2>|<h3>")
	rule2 := regexp.MustCompile("</h1>|</h2>|</h3>")
	rule3 := regexp.MustCompile("<span>|</span>|<div>|</div>")
	s := rule1.ReplaceAllString(*str, "<h4>")
	s = rule2.ReplaceAllString(s, "</h4>")
	s = rule3.ReplaceAllString(s, "")
	return  s
}