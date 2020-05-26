make:
	go1.14.3 build -o bin/scraper cmd/main.go 

build:
	make

run:
	go1.14.3 run cmd/main.go

run-build:
	./bin/scraper
