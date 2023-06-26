.SILENT:

build:
	go build -o ./.bin/app main.go

lint:
	golangci-lint run

docker:
	docker build -t linxdatacenter .
