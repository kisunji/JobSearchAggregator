# Job Search Aggregator

Collection of services that feeds an API with job data that I am interested in.

Started this project to learn golang and to save time while job searching

## Requirements
* Go 1.12^

## How to Run
`go run main.go`

This application will use go http server on port 80 by default. You can set the following environment variables:

`MODE=lambda` if specified, app will use aws-lambda-go's handler

`PORT=8080` set custom port number (only if not running on MODE=lambda)

`CORS=https://yoururlhere.com` set custom origin cor CORS (default is "*")