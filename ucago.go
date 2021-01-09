package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/spf13/viper"
	//	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}

var date_month string = ""

func main() {
	// create new collectors
	collector := colly.NewCollector()
	details_collector := colly.NewCollector()

	// data container
	cal := NewCalendar()

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
		log.Println("Response received", r.StatusCode)
	})
	/*
		details_collector.OnResponse(func(r *colly.Response) {
			log.Println("Details_collector response received", r.StatusCode)
		})
	*/
	reCalDate := regexp.MustCompile(`[0-9][0-9]?/([0-9]{2})?`)

	collector.OnHTML(".ZhCalMonthDay", func(e *colly.HTMLElement) {
		date := e.ChildText("div > a[href]")

		// format the date to dd/mm/yyyy from the zimbra date dd or dd/mm
		if reCalDate.MatchString(date) {
			date_month = strings.Split(date, "/")[1]
			date = fmt.Sprintf("%s", lpad(date, "0", 5))
		} else {
			date = fmt.Sprintf("%s/%s", lpad(date, "0", 2), date_month)
		}

		// for each course of the day
		e.ForEach(".ZhCalMonthAppt", func(ind int, item *colly.HTMLElement) {
			course_link := item.ChildAttr("a[href]", "href")
			course_name := item.ChildText("a[href]")

			course_name = strings.Join(strings.Fields(strings.TrimSpace(course_name)), " ")
			splited := strings.SplitN(course_name, " ", 2) // split HH:MM and NAME from course_name = "HH:MM NAME"

			course := NewCourse(splited[0], splited[1])
			cal.AddCourse(date, course)

			details_collector.Visit(e.Request.AbsoluteURL(course_link))
		})

	})

	reDate := regexp.MustCompile(`.+?,`)
	reTime := regexp.MustCompile(`[0-9][0-9]:[0-9][0-9]`)
	reLink := regexp.MustCompile(`https?://teams.microsoft.com/l/meetup-join/.+?>`)
	details_collector.OnHTML("table.Compose", func(e *colly.HTMLElement) {
		content := e.ChildText("#iframeBody.MsgBody")

		url := strings.TrimSuffix(string(reLink.Find([]byte(content))), ">")

		var start_end []string
		var date string
		header := make(map[string]string)

		e.ForEach(".MsgHdr tbody tr", func(ind int, item *colly.HTMLElement) {
			hdr_name := item.ChildText(".MsgHdrName")
			hdr_value := item.ChildText(".MsgHdrValue")

			hdr_name = strings.Join(strings.Fields(strings.Trim(strings.TrimSpace(hdr_name), " :")), " ")
			hdr_value = strings.Join(strings.Fields(strings.TrimSpace(hdr_value)), " ")

			if hdr_name == "Date" {
				start_end = reTime.FindAllString(hdr_value, -1)
				date = GetDateFormat(reDate.FindString(hdr_value))
			} else if hdr_name == "Participants" && len(hdr_value) > 100 {
				hdr_value = "Long list"
			}
			header[hdr_name] = hdr_value
		})

		cal.CourseList[date][start_end[0]].Link = url
		cal.CourseList[date][start_end[0]].End = start_end[1]
		cal.CourseList[date][start_end[0]].Info = header
	})

	// Before making a request print "Visiting ..."
	collector.OnRequest(func(r *colly.Request) {
		log.Println("Scrapping", r.URL.String())
	})

	// start scraping on  zimbra/h for basic client without js
	collector.Visit("https://mail.uca.fr/zimbra/h/calendar?view=month")

	for k, v := range cal.CourseList {
		fmt.Println("| DATE |:", k)
		for k2, v2 := range v {
			fmt.Printf("\t ---%s---\n", k2)
			v2.Display()
			fmt.Println("---------------------------------------------------------------")
		}
	}
}
