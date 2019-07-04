package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kisunji/jobsearchaggregator/jobservice"
)

var (
	// ErrJobService is thrown when there is an issue unmarshalling the json produced by jobservice
	ErrJobService = errors.New("There was an issue with the jobservice API")
)

//Handler is the AWS Lambda function handler that uses Amazon API Gateway request/response
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	defer timeTrack(time.Now(), "Handler")
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)
	c := make(chan []jobservice.Job)
	go func() { c <- jobservice.AmazonJobs() }()
	go func() { c <- jobservice.LeagueJobs() }()
	var jobArray []jobservice.Job
	for i := 0; i < 2; i++ {
		result := <-c
		jobArray = append(jobArray, result...)
	}
	log.Printf("Jobs found: %d", len(jobArray))
	bytes, err := json.Marshal(jobArray)
	if err != nil {
		return events.APIGatewayProxyResponse{}, ErrJobService
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(bytes),
	}, nil
}

func main() {
	lambda.Start(Handler)
}

// timeTrack measures time to execute
// Credits: https://blog.stathat.com/2012/10/10/time_any_function_in_go.html
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
