## Have neccessary docker image running by using the command
docker-compose up -d
## Application 1
To start application 1 you can run the following to get help 

* go run ./cmd/api/ -help

To start application 1
* go run .cmd/api/
To run unit test
* go test -v ./cmd/api/

## Application 2

To start application 2 you can run the following to get help 

* go run ./cmd/messageProcessor/main.go

To start application 2
* go run .cmd/api/
To run unit test
* go test -v ./cmd/api/

## Application 3
To start application 3 you can run the following to get help 

* go run ./cmd/reportingApi/ -help

To start application 3
* go run ./cmd/reportingApi/
To run unit test
* go test -v ./cmd/reportingApi

## NOTE
For all commands to run properly, run them in the root directory