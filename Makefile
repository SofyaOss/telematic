build:
	sudo docker compose build

run:
	sudo docker compose up app

go-test:
	go test ./...

lint:
	golint ./...
