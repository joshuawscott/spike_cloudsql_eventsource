package watcher

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

// PostgresWatcher is a Watcher specific to postgresql.
type PostgresWatcher struct {
	// Name of the audit table
	AuditTable string
	// ConnectionString is the URL used to connect.
	ConnectionString string
	// connection is the connection to the postgres database
	connection *pgx.Conn
	// Name of the PL/PGSQL Function
	Function string
	// Name of the channel events fire on inside postgres.
	PubSub string
	// Prefix of the Trigger (actual triggers are named <Trigger>_<table name>)
	Trigger string
}

// NewPostgresWatcher returns a struct that allows creating the audit infrastructure for the watcher.
func NewPostgresWatcher() (PostgresWatcher, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PGURL"))
	return PostgresWatcher{
		AuditTable:       "watcher_audit",
		Function:         "watcher_audit_function",
		Trigger:          "watcher_audit_trigger",
		PubSub:           "watcher_audit_channel",
		ConnectionString: os.Getenv("PGURL"),
		connection:       conn,
	}, err
}

// CreateAuditTable implements Watcher.
func (w PostgresWatcher) CreateAuditTable() error {
	query := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
		inserted_at timestamp with time zone DEFAULT now(),
		table_name text,
		old_data jsonb,
		new_data jsonb
	)
	`, w.AuditTable)

	_, err := w.connection.Exec(context.Background(), query)
	return err
}

// CreateFunction creates a function that the trigger will use (required implements Watcher.
func (w PostgresWatcher) CreateFunction() error {
	query := fmt.Sprint(`
	CREATE OR REPLACE FUNCTION `, w.Function, `() RETURNS trigger AS $audit_trail_watcher$
	DECLARE
		inserted `, w.AuditTable, `%ROWTYPE;
		ins_id uuid;
	BEGIN
		IF (TG_OP = 'UPDATE') THEN
			INSERT INTO `, w.AuditTable, `(table_name, old_data, new_data) VALUES (TG_TABLE_NAME, row_to_json(OLD.*), row_to_json(NEW.*)) RETURNING id into ins_id;
		ELSIF (TG_OP = 'INSERT') THEN
			INSERT INTO `, w.AuditTable, `(table_name, old_data, new_data) VALUES (TG_TABLE_NAME, NULL, row_to_json(NEW.*)) RETURNING id into ins_id;
		ELSIF (TG_OP = 'DELETE') THEN
			INSERT INTO `, w.AuditTable, `(table_name, old_data, new_data) VALUES (TG_TABLE_NAME, row_to_json(OLD.*), NULL) RETURNING id into ins_id;
		END IF;
		PERFORM pg_notify('`, w.PubSub, `'::text, ins_id::text);
		RETURN NULL;
	END;
	$audit_trail_watcher$ LANGUAGE plpgsql;
	`)
	fmt.Println(query)

	_, err := w.connection.Exec(context.Background(), query)
	return err
}

// CreateTrigger Implements Watcher
func (w PostgresWatcher) CreateTrigger(table string) error {
	query := fmt.Sprint("DROP TRIGGER IF EXISTS ", w.Trigger, " ON ", table, "; CREATE TRIGGER ", w.Trigger, " AFTER INSERT OR UPDATE OR DELETE ON ", table, " FOR EACH ROW EXECUTE PROCEDURE ", w.Function, "();")

	_, err := w.connection.Exec(context.Background(), query)
	return err
}

// Watch begins watching for changes and writing them to the returned channel
func (w PostgresWatcher) Watch() chan Notification {
	w.connection.Exec(context.Background(), fmt.Sprint("LISTEN ", w.PubSub))
	channel := make(chan Notification)
	go func() {
		for {
			notification, err := w.connection.WaitForNotification(context.Background())
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error waiting for notification:", err)
				os.Exit(1)
			}

			fmt.Println("PID:", notification.PID, "Channel:", notification.Channel, "Payload:", notification.Payload)

			channel <- Notification{}
		}
	}()

	return channel
}
