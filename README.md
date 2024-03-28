[![Github CI/CD](https://img.shields.io/github/workflow/status/IskenT/MultiGameServices/Go)](https://github.com/IskenT/MultiGameServices/actions)
[![Go Report](https://goreportcard.com/badge/github.com/IskenT/MultiGameServices)](https://goreportcard.com/report/github.com/IskenT/MultiGameServices)
![Repository Top Language](https://img.shields.io/github/languages/top/IskenT/MultiGameServices)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/IskenT/MultiGameServices)
![Github Repository Size](https://img.shields.io/github/repo-size/IskenT/MultiGameServices)
![Github Open Issues](https://img.shields.io/github/issues/IskenT/MultiGameServices)
![License](https://img.shields.io/badge/license-MIT-green)
![GitHub last commit](https://img.shields.io/github/last-commit/IskenT/MultiGameServices)
![GitHub contributors](https://img.shields.io/github/contributors/IskenT/MultiGameServices)
![Simply the best ;)](https://img.shields.io/badge/simply-the%20best%20%3B%29-orange)

# MultiGameServices

## About project
- Three servers are implemented in the project: http, gRPC and WebSocket;
- For data storage in-memory-cache was used (data type `uuid` (user id) was chosen as the key, and `int` was chosen to store the amount of money for optimal storage (data type `float` can be implemented for input and multiplied by 100));
- To limit the storage of data in the cache, the time must be set `ttl`;
- The project is logged (writes to the console by default);
- The project is dockerized and a multistage docker file is used;

##  Startup instructions
Docker-compose is required to run the project.

`docker-compose up --build --remove-orphans` or `make run`

Proto files are located in the proto directory, to start run `make run`

## API
- REST Endpoint = `http://localhost:8080/`
- gRPC Endpoints = `0.0.0.0:9090` (In Postman, use gRPC format requests with Proto file loading)
- WebSocket = `http://localhost:8081/`

## A picture is worth a thousand words
<img src="./images/runner.PNG">

<img src="./images/socket.PNG">