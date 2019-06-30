package jobservice

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	leagueURL   = "https://league.com/ca/careers-at-league/jobs/"
	wordFilters = "Director,Senior"
)

// LeagueJobs does something
func LeagueJobs() {
	doc := getDocument(leagueURL)
	jobs := getJobs(doc)
	filteredJobs := filterJobs(jobs, strings.Split(wordFilters, ","))
	listOfURLs := extractURLs(filteredJobs)
	for _, v := range listOfURLs {
		getJobPosting(v)
	}
}

func getDocument(URL string) *goquery.Document {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func getJobs(doc *goquery.Document) *goquery.Selection {
	results := doc.Find(".job-openings__container__jobs__job").FilterFunction(
		func(i int, s *goquery.Selection) bool {
			for _, v := range s.Children().Nodes {
				if v.Data == "h4" && v.FirstChild.Data == "Engineering" {
					return true
				}
			}
			return false
		})
	return results.Find(".job-openings__container__jobs--list").Children()
}

func filterJobs(jobs *goquery.Selection, wordFilters []string) *goquery.Selection {
	return jobs.FilterFunction(
		func(i int, s *goquery.Selection) bool {
			for _, v := range wordFilters {
				if strings.Contains(s.Text(), v) {
					return false
				}
			}
			return true
		})
}

func extractURLs(filteredJobs *goquery.Selection) []string {
	return filteredJobs.Map(
		func(i int, s *goquery.Selection) string {
			link, ok := s.Attr("href")
			if !ok {
				log.Fatal("No href found")
			}
			return link
		})
}

func getJobPosting(url string) {
	doc := getDocument(url)
	titleSelection := doc.Find("meta").FilterFunction(
		func(i int, s *goquery.Selection) bool {
			v, _ := s.Attr("property")
			return v == "og:title"
		})
	title, ok := titleSelection.Attr("content")
	if !ok {
		log.Print("No title found")
		return
	}
	log.Print(title)
	descriptionSelection := doc.Find("meta").FilterFunction(
		func(i int, s *goquery.Selection) bool {
			v, _ := s.Attr("property")
			return v == "og:description"
		})
	description, ok := descriptionSelection.Attr("content")
	if !ok {
		log.Print("No description found")
	}
	log.Print(description)
	requirementSelection := doc.Find("ul.posting-requirements.plain-list > ul > li")
	var requirementArr []string
	for _, n := range requirementSelection.Nodes {
		requirementArr = append(requirementArr, n.FirstChild.Data)
	}
	requirements := strings.Join(requirementArr, "<br/>")
	log.Print(requirements)

}
