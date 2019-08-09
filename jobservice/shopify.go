package jobservice

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ShopifySearch holds the base search URL and an array of words to filter
type ShopifySearch struct {
	SearchURL   string
	BaseURL     string
	wordFilters []string
}

// NewShopifySearch returns the default ShopifySearch
func NewShopifySearch() *ShopifySearch {
	return &ShopifySearch{
		SearchURL:   "https://www.shopify.ca/careers/search?specialties%5B%5D=1&locations%5B%5D=2&keywords=&sort=",
		BaseURL:     "https://www.shopify.ca",
		wordFilters: []string{"Director", "Senior", "Manager", "Lead", "Welcome"},
	}
}

// Jobs calls Shopify's careers page and parses results based on specific css selectors
func (s *ShopifySearch) Jobs() []Job {
	var jobArray []Job
	doc := s.getDocument(s.SearchURL)
	jobs := s.getJobs(doc)
	filteredJobs := s.filterJobs(jobs, s.wordFilters)
	listOfURLs := s.extractURLs(filteredJobs)
	for _, v := range listOfURLs {
		jobArray = append(jobArray, s.getJobPosting(v))
	}
	return jobArray
}

func (s *ShopifySearch) getDocument(URL string) *goquery.Document {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func (s *ShopifySearch) getJobs(doc *goquery.Document) *goquery.Selection {
	return doc.Find(".jobs-table__cell").Children()
}

func (s *ShopifySearch) filterJobs(jobs *goquery.Selection, wordFilters []string) *goquery.Selection {
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

func (s *ShopifySearch) extractURLs(filteredJobs *goquery.Selection) []string {
	return filteredJobs.Map(
		func(i int, s *goquery.Selection) string {
			link, ok := s.Attr("href")
			if !ok {
				log.Fatal("No href found")
			}
			return link
		})
}

func (s *ShopifySearch) getJobPosting(url string) Job {
	doc := s.getDocument(s.BaseURL + url)

	title := s.getTitle(doc)
	description := s.getDescription(doc)
	requirementArr := s.getRequirements(doc)

	return Job{
		Company:        "Shopify",
		Title:          title,
		Description:    description,
		Qualifications: requirementArr,
		URL:            s.BaseURL + url,
	}
}

func (s *ShopifySearch) getTitle(doc *goquery.Document) string {
	return s.getFieldFromMeta(doc, "og:title")
}

func (s *ShopifySearch) getDescription(doc *goquery.Document) string {
	return s.getFieldFromMeta(doc, "og:description")
}

func (s *ShopifySearch) getFieldFromMeta(doc *goquery.Document, fieldName string) string {
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

func (s *ShopifySearch) getRequirements(doc *goquery.Document) []string {
	var requirementArr []string
	requirementSelection := doc.Find("ul.job-posting__list > li")
	for _, n := range requirementSelection.Nodes {
		requirementArr = append(requirementArr, n.FirstChild.Data)
	}
	return requirementArr
}
