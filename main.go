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

// Handler is the AWS Lambda function handler that uses Amazon API Gateway request/response
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	defer timeTrack(time.Now(), "Handler")
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)
	var jobArray []jobservice.Job
	jobArray = append(jobArray, jobservice.AmazonJobs()...)
	jobArray = append(jobArray, jobservice.LeagueJobs()...)
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
// Taken from https://blog.stathat.com/2012/10/10/time_any_function_in_go.html
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
