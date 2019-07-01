package main

import (
	"log"

	"github.com/kisunji/jobsearchaggregator/jobservice"
)

func main() {
	jobservice.AmazonJobs()
	log.Print(jobservice.LeagueJobs())
}
