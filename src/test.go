package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const amazonURL = "https://www.amazon.jobs/en/search.json?base_query=&category[]=software-development&job_function_id[]=job_function_corporate_80rdb4&=&normalized_location[]=Toronto,+Ontario,+CAN&offset=0&query_options=&radius=24km&region=&result_limit=200&sort=recent"

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

// Filter job slice based on a predicate func
func Filter(vs []Job, f func(Job) bool) []Job {
	vsf := make([]Job, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func isSuitable(job Job) bool {
	// Positions containing these words are generally not suitable
	if strings.Contains(job.Title, "Manager") ||
		strings.Contains(job.Title, "Senior") ||
		strings.Contains(job.Title, "Sr") ||
		strings.Contains(job.Title, "II") {
		return false
	}

	// If there is mention of numbers of years, keep it to 3 or less
	re := regexp.MustCompile(`[4-9]\+? year`)
	if re.MatchString(job.Qualifications) {
		return false
	}
	return true
}

func isRecent(job Job) bool {
	// Make sure job was updated within last 2 months
	if strings.Contains(job.TimeSinceLastUpdated, "month") {
		re := regexp.MustCompile(`[0-9]+`)
		monthString := re.FindString(job.TimeSinceLastUpdated)
		monthValue, err := strconv.Atoi(monthString)
		if err != nil {
			log.Fatal(err)
		}
		return monthValue <= 2
	}
	// If posting contains the word "year", ignore it
	if strings.Contains(job.TimeSinceLastUpdated, "year") {
		return false
	}
	return true
}

func main() {
	fmt.Println("Hello World!")
	resp, err := http.Get(amazonURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if !json.Valid(body) {
		log.Fatal("Not a valid Json", err)
	}
	var jobList JobList
	err = json.Unmarshal(body, &jobList)
	log.Printf("Number of jobs detected: %d", len(jobList.Jobs))

	suitableJobs := Filter(jobList.Jobs, isSuitable)
	suitableJobs = Filter(suitableJobs, isRecent)
	log.Printf("Number of suitable jobs detected: %d", len(suitableJobs))

	for _, v := range suitableJobs {
		log.Println(v)
	}
}
