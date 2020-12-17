# Spike of CloudSQL Eventsourcing

## Setup

* Create .env file like:
```
PGPASSWORD=<PASSWORD HERE>
PGURL=postgres://postgres:<PASSWORD HERE>@localhost:5432/postgres
```
inserting a password in `<PASSWORD HERE>`

## Current status / goals
* [x] Creates trigger to send notifications.
* [x] Reads notification.
* [ ] Writes data into downstream.

## Trying it out

Check out the makefile for more detail.

1. Setup - brings up docker compose, compiles, and creates tables in the docker postgres.
```
make setup
```

2. Open a psql CLI session inside docker:
```
make psql
```

3. Run the watcher in another terminal:
```
# source the .env environment variables
export $(cat .env | xargs)
# Run the watcher
./spike_cloudsql_eventsource -watch
```

4. In the psql CLI session, insert into one of the test tables (`table_incr` or `table_uuid` - see below for fields).

Example query:
```
insert into table_incr (bool_field, text_field, bigint_field) values (true, 'hello', 100);
```

You should see the watcher window print something like:
```
PID: 100 Channel: watcher_audit_channel Payload: e0753132-c2b2-4167-9fab-4aa70d98f141
```

The payload is the UUID of the record inserted into the watcher_audit table.

## Hacking on it

Running `make setup` again should update the trigger function definitions, but
if you change more than that, you should redo the setup with
```
make clean setup
```

This will kill the database and re-compile and re-run the setup steps.

### Test Table layouut

```
postgres=# \d table_incr
                                 Table "public.table_incr"
    Column    |  Type   |                             Modifiers
--------------+---------+-------------------------------------------------------------------
 serial_field | bigint  | not null default nextval('table_incr_serial_field_seq'::regclass)
 bool_field   | boolean |
 text_field   | text    |
 bigint_field | bigint  |
Indexes:
    "table_incr_pkey" PRIMARY KEY, btree (serial_field)
Triggers:
    watcher_audit_trigger AFTER INSERT OR DELETE OR UPDATE ON table_incr FOR EACH ROW EXECUTE PROCEDURE watcher_audit_function()

postgres=# \d table_uuid
                  Table "public.table_uuid"
    Column    |  Type   |              Modifiers
--------------+---------+-------------------------------------
 uuid_field   | uuid    | not null default uuid_generate_v4()
 bool_field   | boolean |
 text_field   | text    |
 bigint_field | bigint  |
Indexes:
    "table_uuid_pkey" PRIMARY KEY, btree (uuid_field)
Triggers:
    watcher_audit_trigger AFTER INSERT OR DELETE OR UPDATE ON table_uuid FOR EACH ROW EXECUTE PROCEDURE watcher_audit_function()
```
