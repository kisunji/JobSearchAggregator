package leagueservice

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const leagueURL = "https://league.com/ca/careers-at-league/jobs/"

// GetSearchResults does something
func GetSearchResults() {
	doc := getDocNode(leagueURL)
	results := doc.Find(".job-openings__container__jobs__job").
		FilterFunction(
			func(i int, s *goquery.Selection) bool {
				for _, v := range s.Children().Nodes {
					if v.Data == "h4" && v.FirstChild.Data == "Engineering" {
						return true
					}
				}
				return false
			})
	jobsContainer := results.Find(".job-openings__container__jobs--list")
	filteredJobs := jobsContainer.Children().FilterFunction(func(i int, s *goquery.Selection) bool {
		return !strings.Contains(s.Text(), "Director") && !strings.Contains(s.Text(), "Senior")
	})
	listOfURLs := filteredJobs.Map(func(i int, s *goquery.Selection) string {
		link, _ := s.Attr("href")
		return link
	})
	log.Print(listOfURLs)
	for _, v := range listOfURLs {
		getJobPosting(v)
	}
}

func getDocNode(URL string) *goquery.Document {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func getJobPosting(url string) {
	doc := getDocNode(url)
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
