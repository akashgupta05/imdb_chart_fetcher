package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

type ImdbChart struct {
	Title            string `json:"title"`
	MovieReleaseYear string `json:"movie_release_year"`
	ImdbRating       string `json:"imdb_rating"`
	Summary          string `json:"summary"`
	Duration         string `json:"duration"`
	Genre            string `json:"genre"`
}

type ImdbCharts struct {
	Charts []ImdbChart
	Error  error
}

func NewImdbCharts() *ImdbCharts {
	return &ImdbCharts{}
}

func (imdbCharts *ImdbCharts) fetch(url string, count int) {
	c := colly.NewCollector(
		colly.Async(true),
		colly.MaxDepth(2),
	)

	c.Limit(&colly.LimitRule{
		Parallelism: 10,
	})

	counter := 0
	c.OnHTML("td.posterColumn > a", func(e *colly.HTMLElement) {
		err := e.Request.Visit(e.Attr("href"))
		if err != nil {
			counter++
			fmt.Printf("error while request to visit url: %v", err.Error())
		}

		if e.Index-counter > count {
			return
		}
	})

	c.OnHTML("#title-overview-widget", func(element *colly.HTMLElement) {
		title := element.ChildText(".titleBar h1")
		title = regexp.MustCompile(`\(\d{4}\)`).ReplaceAllString(title, "")
		title = strings.TrimLeft(title, " ")

		year := element.ChildText("#titleYear")
		year = strings.ReplaceAll(year, "(", "")
		year = strings.ReplaceAll(year, ")", "")

		var movie = ImdbChart{
			Title:            title,
			MovieReleaseYear: year,
			ImdbRating:       element.ChildText("div.ratingValue > strong > span"),
			Summary:          element.ChildText(".summary_text"),
			Duration:         element.ChildText("time"),
			Genre:            element.ChildText("div.subtext > a:nth-child(4)"),
		}

		imdbCharts.Charts = append(imdbCharts.Charts, movie)
	})

	c.Visit(url)
	c.Wait()
}
