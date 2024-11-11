SERVICE_NAME := s3-downloader

init:
	make clean

tidy:
	go mod tidy

build:
	go build -o $(SERVICE_NAME) main.go

vendor:
	make clean
	go mod tidy
	go mod vendor

clean-dbs:
	rm -rf dbs

clean-build:
	rm -f $(SERVICE_NAME)

clean:
	make clean-dbs
	make clean-build
	rm -rf vendor

compose-up:
	docker compose up -d

compose-down:
	docker compose down
