package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
	//	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}

func main() {
	// create a new collector
	collector := colly.NewCollector()
	details_collector := colly.NewCollector()

	// change this line for command input instead of .env
	body := "uusername=" + viper.GetString("USERNAME") + "&password=" + viper.GetString("PASSWORD") + "&execution=" + viper.GetString("EXECUTION") + "&_eventId=submit&submit=LOGIN"

	// authentication
	err := collector.PostRaw("https://ent.uca.fr/cas/login?service=https%3A%2F%2Fmail.uca.fr%2Fzimbra%2Fpublic%2Fpreauthuca.jsp", []byte(body))
	if err != nil {
		log.Fatal(err)
	}
	err = details_collector.PostRaw("https://ent.uca.fr/cas/login?service=https%3A%2F%2Fmail.uca.fr%2Fzimbra%2Fpublic%2Fpreauthuca.jsp", []byte(body))
	if err != nil {
		log.Fatal(err)
	}

	// attach callbacks
	collector.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		//r.Save("./body.html")
	})

	details_collector.OnResponse(func(r *colly.Response) {
		log.Println("DC response received", r.StatusCode)
		//r.Save("./body.html")
	})

	// On every a element which has href attribute call callback
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		//link := e.Attr("href")
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		//c.Visit(e.Request.AbsoluteURL(link))
	})
	collector.OnHTML("#mess_list_tbody tr", func(e *colly.HTMLElement) {
		// data := e.ChildText("td[colspan='3']")
		content_link := e.ChildAttr("td > a[href]", "href")

		// data = strings.ReplaceAll(data, " \t", "|")
		// val := strings.Split(strings.Join(strings.Fields(strings.TrimSpace(data)), " "), "|")
		// fmt.Println(val, "link:", content_link)
		details_collector.Visit(e.Request.AbsoluteURL(content_link))

	})
	details_collector.OnHTML(".ZhAppContent2", func(e *colly.HTMLElement) {
		data := e.ChildTexts(".MsgHdr table table tr")
		date := e.ChildText("td[class='MsgHdrSent'][align='right']")
		links := e.ChildAttrs("#iframeBody a[href][class='zUrl']", "href")
		//		content := e.ChildText("#iframeBody")

		re := regexp.MustCompile(`https://teams.microsoft.com/l/meetup-join/`)
		ms_links := make([]string, 0)

		for _, l := range links {
			if re.MatchString(l) {
				ms_links = append(ms_links, l)
			}
		}

		header := make(map[string]string)
		cpt := 0

		reDate := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
		for _, val := range data {
			vals := strings.Split(strings.Join(strings.Fields(strings.TrimSpace(val)), " "), ":")
			if len(vals) < 2 || reDate.MatchString(val) {
				header["other_"+strconv.Itoa(cpt)] = strings.Join(strings.Fields(strings.TrimSpace(val)), " ")
				cpt += 1
			} else {
				header[vals[0]] = vals[1]
			}
		}

		if len(ms_links) > 0 {
			fmt.Println(ms_links)
			fmt.Println(date)
			fmt.Println("\tHeader:")
			for k, v := range header {
				fmt.Println("\t\t", k, "->", v)
			}
		}
		/*		mails := &Mail{
					Content:  val[0],
					Links : ms_links
					Header: header,
				}
		*/
	})

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		// log.Println("Visiting", r.URL.String())
	})

	/*
		details_collector.OnRequest(func(r *colly.Request) {
			log.Println("Visiting", r.URL.String())
		})
	*/

	// start scraping
	collector.Visit("https://mail.uca.fr/zimbra/h/") // zimbra/h for basic client without js
}
