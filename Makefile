# Import the .env file
ifneq (,$(wildcard ./.env))
include .env
	export
endif

test:
	go test -ldflags "-X spike_cloudsql_eventsource/pkg/watcher.postgresPassword=${POSTGRES_PASSWORD}" ./...
