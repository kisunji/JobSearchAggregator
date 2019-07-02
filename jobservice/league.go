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
func LeagueJobs() []Job {
	var jobArray []Job
	doc := getDocument(leagueURL)
	jobs := getJobs(doc)
	filteredJobs := filterJobs(jobs, strings.Split(wordFilters, ","))
	listOfURLs := extractURLs(filteredJobs)
	for _, v := range listOfURLs {
		jobArray = append(jobArray, getJobPosting(v))
	}
	return jobArray
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

func getJobPosting(url string) Job {
	doc := getDocument(url)

	title := getTitle(doc)
	description := getDescription(doc)
	requirementArr := getRequirements(doc)
	requirements := strings.Join(requirementArr, "<br/>")

	return Job{
		Company:        "League Inc.",
		Title:          title,
		Description:    description,
		Qualifications: requirements,
		URL:            url}
}

func getTitle(doc *goquery.Document) string {
	return getFieldFromMeta(doc, "og:title")
}

func getDescription(doc *goquery.Document) string {
	return getFieldFromMeta(doc, "og:description")
}

func getFieldFromMeta(doc *goquery.Document, fieldName string) string {
	selection := doc.Find("meta").FilterFunction(
		func(i int, s *goquery.Selection) bool {
			v, _ := s.Attr("property")
			return v == fieldName
		})
	field, ok := selection.Attr("content")
	if !ok {
		log.Printf("Field: %s not found", fieldName)
		return "Field not found"
	}
	return field
}

func getRequirements(doc *goquery.Document) []string {
	var requirementArr []string
	requirementSelection := doc.Find("ul.posting-requirements.plain-list > ul > li")
	for _, n := range requirementSelection.Nodes {
		requirementArr = append(requirementArr, n.FirstChild.Data)
	}
	return requirementArr
}
