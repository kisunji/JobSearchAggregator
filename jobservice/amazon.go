package jobservice

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	amazonURL     = "https://www.amazon.jobs/en/search.json?base_query=&category[]=software-development&job_function_id[]=job_function_corporate_80rdb4&=&normalized_location[]=Toronto,+Ontario,+CAN&offset=0&query_options=&radius=24km&region=&result_limit=200&sort=recent"
	amazonBaseURL = "https://www.amazon.jobs"
)

// amazonJob holds a subset of fields I care about
type amazonJob struct {
	Title                   string
	Description             string `json:"description_short"`
	Qualifications          string `json:"basic_qualifications"`
	PreferredQualifications string `json:"preferred_qualifications"`
	DatePosted              string `json:"posted_date"`
	Category                string `json:"business_category"`
	URL                     string `json:"job_path"`
	TimeSinceLastUpdated    string `json:"updated_time"`
}

// amazonJobList represents the highest level struct returned by Amazon API
type amazonJobList struct {
	Jobs []amazonJob
}

// AmazonJobs calls Amazon's job search API and applies custom filters to show only relevant job postings
func AmazonJobs() []Job {
	var jobArray []Job
	responseBody := callAPI(amazonURL)
	jobList := convertToJSONList(responseBody)

	suitableJobs := filter(jobList.Jobs, isRecent, isSuitable)

	for _, v := range suitableJobs {
		qualifications := processQualifications(v.Qualifications)
		prefQualifications := processQualifications(v.PreferredQualifications)

		jobArray = append(jobArray, Job{
			Company:                 "Amazon",
			Title:                   v.Title,
			Qualifications:          qualifications,
			PreferredQualifications: prefQualifications,
			Description:             v.Description,
			URL:                     amazonBaseURL + v.URL,
		})
	}
	return jobArray
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

func convertToJSONList(bytes []byte) amazonJobList {
	if !json.Valid(bytes) {
		log.Fatal("Not a valid Json")
	}
	var jobList amazonJobList
	err := json.Unmarshal(bytes, &jobList)
	if err != nil {
		log.Fatal(err)
	}
	return jobList
}

// filters based on any number of predicates
// most restrictive filter (likely to false) should be passed first
func filter(vs []amazonJob, fs ...func(amazonJob) bool) []amazonJob {
	vsf := make([]amazonJob, 0)
OUTER:
	for _, v := range vs {
		for _, f := range fs {
			if !f(v) {
				continue OUTER
			}
		}
		vsf = append(vsf, v)
	}
	return vsf
}

func isSuitable(j amazonJob) bool {
	// Positions containing these words are generally not suitable
	if strings.Contains(j.Title, "Manager") ||
		strings.Contains(j.Title, "Senior") ||
		strings.Contains(j.Title, "Sr") ||
		strings.Contains(j.Title, "II") {
		return false
	}

	// If there is mention of numbers of years, keep it to 3 or less
	re := regexp.MustCompile(`[4-9]\+? year`)
	if re.MatchString(j.Qualifications) {
		return false
	}
	return true
}

func isRecent(j amazonJob) bool {
	// Make sure job was updated within last 2 months
	if strings.Contains(j.TimeSinceLastUpdated, "month") {
		re := regexp.MustCompile(`[0-9]+`)
		monthString := re.FindString(j.TimeSinceLastUpdated)
		monthValue, err := strconv.Atoi(monthString)
		if err != nil {
			log.Fatal(err)
		}
		return monthValue <= 1
	}
	// If posting contains the word "year", ignore it
	if strings.Contains(j.TimeSinceLastUpdated, "year") {
		return false
	}
	return true
}

func processQualifications(rawQualifications string) []string {
	var result []string

	// Trim unwanted bullet rune from raw string
	q := strings.Replace(rawQualifications, "•", "", -1)
	// Since amazon's api returns strings laced with <br/>, use the tags to split into array
	qs := strings.Split(q, "<br/>")

	for _, q := range qs {
		q = strings.TrimSpace(q)
		// trimmed string should not:
		// be empty, be a hashtag, be disclaimer
		if !(len(q) == 0 || strings.HasPrefix(q, "#") || strings.HasPrefix(q, "Amazon is an")) {
			result = append(result, q)
		}
	}
	return result
}
