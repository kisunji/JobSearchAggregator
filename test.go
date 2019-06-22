package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const amazonURL = "https://www.amazon.jobs/en/search.json?base_query=&category[]=software-development&job_function_id[]=job_function_corporate_80rdb4&=&normalized_location[]=Toronto,+Ontario,+CAN&offset=0&query_options=&radius=24km&region=&result_limit=200&sort=relevant"

// Job holds a subset of fields I care about
type Job struct {
	Title                   string
	Qualifications          string `json:"basic_qualifications"`
	PreferredQualifications string `json:"preferred_qualifications"`
	DatePosted              string `json:"posted_date"`
	Category                string `json:"business_category"`
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
	if strings.Contains(job.Title, "Manager") || strings.Contains(job.Title, "Senior") {
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
	log.Printf("Number of suitable jobs detected: %d", len(suitableJobs))
}
