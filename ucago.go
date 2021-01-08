package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/url"
	"strings"
)

func createFormReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}

func main() {
	// create a new collector
	c := colly.NewCollector()

	body := "uusername=" + viper.GetString("USERNAME") + "&password=" + viper.GetString("PASSWORD") + "&execution=" + viper.GetString("EXECUTION") + "&_eventId=submit&submit=LOGIN"

	// authenticate
	err := c.PostRaw("https://ent.uca.fr/cas/login?service=https%3A%2F%2Fmail.uca.fr%2Fzimbra%2Fpublic%2Fpreauthuca.jsp", []byte(body))
	if err != nil {
		log.Fatal(err)
	}

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		//		r.Save("./body.html")
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		//link := e.Attr("href")
		// Print link
		//		log.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		//Only those links are visited which are in AllowedDomains
		//		c.Visit(e.Request.AbsoluteURL(link))
	})
	c.OnHTML("#mess_list_tbody tr", func(e *colly.HTMLElement) {
		data := e.ChildText("td[colspan='3']")
		//		myByte := []byte{194, 160}
		data = strings.ReplaceAll(data, string([]byte{194, 160}), "")
		data = strings.ReplaceAll(data, " ", "")
		data = strings.ReplaceAll(data, "\n", "")
		//		title = strings.ReplaceAll(title, "\t", "-")
		vals := strings.Split(data, "\t")
		fmt.Println("source:", vals[0], "content:", vals[1])
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	// start scraping
	c.Visit("https://ent.uca.fr/")
	c.Visit("https://mail.uca.fr/zimbra/h/")
}
