package leagueservice

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const leagueURL = "https://league.com/ca/careers-at-league/jobs/"

// Job holds a subset of fields I care about
type Job struct {
	Title                   string
	Qualifications          string `json:"basic_qualifications"`
	PreferredQualifications string `json:"preferred_qualifications"`
	DatePosted              string `json:"posted_date"`
	Category                string `json:"business_category"`
	URL                     string `json:"job_path"`
	TimeSinceLastUpdated    string `json:"updated_time"`
}

// JobList represents the highest level struct returned by Amazon API
type JobList struct {
	Jobs []Job
}

// GetSearchResults does something
func GetSearchResults() {
	doc := getDocNode(leagueURL)
	results := doc.Find(".job-openings__container__jobs__job").
		FilterFunction(func(i int, s *goquery.Selection) bool {
			for _, v := range s.Children().Nodes {
				if v.Data == "h4" && v.FirstChild.Data == "Engineering" {
					return true
				}
			}
			return false
		})
	log.Print(len(results.Nodes))
	jobsContainer := results.Find(".job-openings__container__jobs--list")
	filteredJobs := jobsContainer.Children().FilterFunction(func(i int, s *goquery.Selection) bool {
		return !strings.Contains(s.Text(), "Director") && !strings.Contains(s.Text(), "Senior")
	})
	filteredJobs.Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		log.Print(link)
	})
}

func getDocNode(URL string) *goquery.Document {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func callAPI(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}
