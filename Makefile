.SILENT:

build:
	sudo docker compose build

run:
	sudo docker compose up app

test:
	go test test/integration_test.go

lint:
	golint ./...
