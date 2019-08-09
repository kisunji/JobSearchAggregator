package jobservice

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// LeagueSearch holds the base search URL and an array of words to filter
type LeagueSearch struct {
	URL         string
	wordFilters []string
}

// NewLeagueSearch returns the default LeagueSearch
func NewLeagueSearch() *LeagueSearch {
	return &LeagueSearch{
		URL:         "https://league.com/ca/careers-at-league/jobs/",
		wordFilters: []string{"Director", "Senior", "Manager"},
	}
}

// Jobs calls League's careers page and parses results based on specific css selectors
func (l *LeagueSearch) Jobs() []Job {
	var jobArray []Job
	doc := l.getDocument(l.URL)
	jobs := l.getJobs(doc)
	filteredJobs := l.filterJobs(jobs, l.wordFilters)
	listOfURLs := l.extractURLs(filteredJobs)
	for _, v := range listOfURLs {
		jobArray = append(jobArray, l.getJobPosting(v))
	}
	return jobArray
}

func (l *LeagueSearch) getDocument(URL string) *goquery.Document {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func (l *LeagueSearch) getJobs(doc *goquery.Document) *goquery.Selection {
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

func (l *LeagueSearch) filterJobs(jobs *goquery.Selection, wordFilters []string) *goquery.Selection {
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

func (l *LeagueSearch) extractURLs(filteredJobs *goquery.Selection) []string {
	return filteredJobs.Map(
		func(i int, s *goquery.Selection) string {
			link, ok := s.Attr("href")
			if !ok {
				log.Fatal("No href found")
			}
			return link
		})
}

func (l *LeagueSearch) getJobPosting(url string) Job {
	doc := l.getDocument(url)

	title := l.getTitle(doc)
	description := l.getDescription(doc)
	requirementArr := l.getRequirements(doc)

	return Job{
		Company:        "League Inc.",
		Title:          title,
		Description:    description,
		Qualifications: requirementArr,
		URL:            url,
	}
}

func (l *LeagueSearch) getTitle(doc *goquery.Document) string {
	return l.getFieldFromMeta(doc, "og:title")
}

func (l *LeagueSearch) getDescription(doc *goquery.Document) string {
	return l.getFieldFromMeta(doc, "og:description")
}

func (l *LeagueSearch) getFieldFromMeta(doc *goquery.Document, fieldName string) string {
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

func (l *LeagueSearch) getRequirements(doc *goquery.Document) []string {
	var requirementArr []string
	requirementSelection := doc.Find("ul.posting-requirements.plain-list > ul > li")
	for _, n := range requirementSelection.Nodes {
		requirementArr = append(requirementArr, n.FirstChild.Data)
	}
	return requirementArr
}
