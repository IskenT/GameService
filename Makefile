run:
	docker-compose up --build --remove-orphans 


lint:
	golangci-lint run