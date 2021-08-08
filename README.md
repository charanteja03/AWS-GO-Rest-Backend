# sf-rest-backend

## Running the project 

1. Copy the .env-example to .env file and provide AWS Credentials 
2. Build the project using
```
go build -o <filename>.exe
```
or just
```
go build
```
If there are errors concerning missing packages import packages with 
`go get <package>`
or import all needed packages
`go get ./...`
3. Run <filename>.exe

## Generating mocks for testing
1. To generate mock mockery is required https://github.com/vektra/mockery
2. Important notice mockery has to be in your environment path for `go generate` to work
3. Run `go generate ./...` on top project directory so mocks are generated at the same level as other packages folders


## Running tests

1. To run test use `go test ./...` this runs all tests
To see what tests have been run while testing add `-v` like this `go test ./... -v`

2. To run tests for directory use 
```
go test <dir>     # For example go test .\healthcheck
```
You can also add `-v` to this command like `go test <dir> -v`
3. To run specyfic test use
```
go test -run ^<testName>$ <dir> # For example go test -run ^TestGetProperMachine$ sfr-backend/machine
```

## Running test coverage 

- To run test coverage use `go test ./... -coverprofile cover.out`
If you want to see report in webbrowser you can use `go tool cover -html=cover.out` command

## Generating Swagger documentation

- To generate swagger:

1. first install go-swagger  `go get -u github.com/go-swagger/go-swagger/cmd/swagger`
2. Generate swagger.yaml  `swagger generate spec -o ./swagger.yaml --scan-models`
3. Serve swagger: `swagger serve -F=swagger swagger.yaml`

IMPORTANT: This just serves the documentation, to enable direct testing from swagger, run your project and change base url on doc.go file to reflect your base path 

