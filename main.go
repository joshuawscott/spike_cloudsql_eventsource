package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"spike_cloudsql_eventsource/pkg/database"
	"spike_cloudsql_eventsource/pkg/watcher"
)

func main() {
	initPtr := flag.Bool("test-tables", false, "Create the test tables")
	watchPtr := flag.Bool("watch", false, "Start the watcher")
	flag.Parse()
	if *initPtr {
		runInit()
		return
	}
	if *watchPtr {
		runWatch()
		return
	}
}

func runInit() {
	conn, err := database.CreateConn()
	if err != nil {
		log.Fatalf("Couldn't connect to postgres: %v\n", err)
	}
	defer conn.Close(context.Background())

	createExtension := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	createExtension2 := `CREATE EXTENSION IF NOT EXISTS "plpgsql";`
	createUUIDTable := `
	CREATE TABLE IF NOT EXISTS table_uuid (
		uuid_field UUID PRIMARY KEY default uuid_generate_v4(),
		bool_field bool,
		text_field text,
		bigint_field bigint
	)
	`
	createIncrTable := `
	CREATE TABLE IF NOT EXISTS table_incr (
		serial_field bigserial PRIMARY KEY,
		bool_field bool,
		text_field text,
		bigint_field bigint
	)
	`
	err = database.CreateTable(conn, createExtension)
	check(err)
	err = database.CreateTable(conn, createExtension2)
	check(err)
	err = database.CreateTable(conn, createUUIDTable)
	check(err)
	err = database.CreateTable(conn, createIncrTable)
	check(err)
}

func runWatch() {
	conn, err := database.CreateConn()
	if err != nil {
		log.Fatalf("Couldn't connect to postgres: %v\n", err)
	}
	defer conn.Close(context.Background())

	pgWatcher, err := watcher.NewPostgresWatcher()
	check(err)
	err = pgWatcher.CreateAuditTable()
	check(err)
	err = pgWatcher.CreateFunction()
	check(err)
	err = pgWatcher.CreateTrigger("table_uuid")
	check(err)
	err = pgWatcher.CreateTrigger("table_incr")
	check(err)

	notifications := pgWatcher.Watch()
	for {
		fmt.Println(<-notifications)
	}
}

// If error then log it and exit.
func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
