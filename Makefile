# Import the .env file
ifneq (,$(wildcard ./.env))
include .env
	export
endif

build:
	go build

setup: docker build
	./spike_cloudsql_eventsource -test-tables

docker:
	docker-compose up -d

clean:
	rm spike_cloudsql_eventsource
	docker-compose down

psql:
	docker-compose exec postgres psql -U postgres
