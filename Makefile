lint:
	golangci-lint run --fix

run:
	docker-compose up fickle

build:
	docker build -f Dockerfile -t artarts36/fickle:0.1.0 .