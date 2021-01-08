package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}

func main() {
	// create a new collector
	c := colly.NewCollector()

	// change this line for command input instead of .env
	body := "uusername=" + viper.GetString("USERNAME") + "&password=" + viper.GetString("PASSWORD") + "&execution=" + viper.GetString("EXECUTION") + "&_eventId=submit&submit=LOGIN"

	// authentication
	err := c.PostRaw("https://ent.uca.fr/cas/login?service=https%3A%2F%2Fmail.uca.fr%2Fzimbra%2Fpublic%2Fpreauthuca.jsp", []byte(body))
	if err != nil {
		log.Fatal(err)
	}

	// attach callbacks
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		//r.Save("./body.html")
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		//link := e.Attr("href")
		//fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		//c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnHTML("#mess_list_tbody tr", func(e *colly.HTMLElement) {
		data := e.ChildText("td[colspan='3']")

		data = strings.ReplaceAll(data, " \t", "|")
		val := strings.Split(strings.Join(strings.Fields(strings.TrimSpace(data)), " "), "|")
		fmt.Println(val)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	// start scraping
	//	c.Visit("https://ent.uca.fr/")
	c.Visit("https://mail.uca.fr/zimbra/h/")
	// zimbra/h for basic client without js
}
