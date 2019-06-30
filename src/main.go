package main

import (
	"github.com/kisunji/jobsearchaggregator/src/amazonservice"
	"github.com/kisunji/jobsearchaggregator/src/leagueservice"
)

func main() {
	leagueservice.GetSearchResults()
	amazonservice.GetSearchResults()
}
