### Overview
This is a simple backend for good old todo service, right now this service can handle login/list/create simple tasks.  
To make it run:
- `go run main.go`
- Import Postman collection from `docs` to check example


### Run unit test with coverate
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
go tool cover -func profile

### Docker build
docker-compose up --build


### TODO list
- Write rernable integration test by postman or jmeter