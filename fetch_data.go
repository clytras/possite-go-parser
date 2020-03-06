package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Event struct {
	Timestamp string   `json:"timestamp"`
	Fulldate  string   `json:"fulldate"`
	Title     string   `json:"title"`
	Subtitle  string   `json:"subtitle"`
	Place     string   `json:"place"`
	Contacts  []string `json:"contacts"`
}

func main() {
	fmt.Printf("Posidonia site parse go test\n\n")

	res, err := http.Get("https://posidonia-events.com/events/conferences-and-seminars/")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	events := doc.Find("body .page-calendar .group ul li.calendar-event .bottom .calendar-event.big")

	if len(events.Nodes) > 0 {
		fmt.Printf("Found %d events; Processing\n", len(events.Nodes))

		data := []Event{}

		events.Each(func(i int, s *goquery.Selection) {
			timestamp := s.Find(".top .when").Text()
			fulldate := s.Find(".top .where .date").Text()
			title := s.Find(".middle .title").Text()
			subtitle := s.Find(".middle .subtitle").Text()
			place := s.Find(".middle .place").Text()

			fmt.Printf(`
Event %d --------------------------------------

Timestamp: %s
Fulldate: %s
Title: %s
Subtitle: %s
Place: %s
`, i+1, timestamp, fulldate, title, subtitle, place)

			contact := s.Find(".middle .contact a")
			contacts := []string{}

			if len(contact.Nodes) > 0 {
				contact.Each(func(ci int, cs *goquery.Selection) {
					email := strings.TrimSpace(cs.Text())

					if len(email) > 0 {
						fmt.Printf("Contact %d: %s\n", ci+1, cs.Text())
						contacts = append(contacts, cs.Text())
					}
				})
			}

			data = append(data, Event{
				timestamp, fulldate, title, subtitle, place, contacts,
			})

			fmt.Println()
		})

		file, _ := json.MarshalIndent(data, "", " ")
		_ = ioutil.WriteFile("pos-data.json", file, 0644)
	} else {
		fmt.Printf("No nodes found\n")
	}

}
