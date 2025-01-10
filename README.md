# NRTMv4 Client

## Introduction

nrtm4client is a tool for communicating with an [NRTMv4 server](https://github.com/mxsasha/nrtmv4) (GROW).
It retrieves IRR data from NTRM mirror servers and stores them in a database (only PostgreSQL
currently). History is maintained and can be queried.

## Usage

Before you can run the client, follow the Quick set up below, then come back to this section.

## Set up environment variables

- NRTM4_FILE_PATH An empty directory where NRTMv4 snapshot and delta files will be stored.
- PG_DATABASE_URL Connection string to PostgreSQL database.

## Running nrtm4client

Create a directory, e.g. `$HOME/nrtm4/RIPE` to store downloaded files,
then copy file [./scripts/env.dev.example.conf] to ./scripts/env.dev.conf, and change the variables
to your system, for example:

    PG_DATABASE_URL=postgres://nrtm4:nrtm4@localhost:5432/nrtm4?sslmode=disable
    NRTM4_FILE_PATH=/tmp/RIPE

Now run the `run*.sh` script in the [./scripts](./scripts) dir like so:

    $ ./scripts/run.sh connect --url <url> # Be patient, snapshot files tend to be on the large side
    $ ./scripts/run.sh list

Command line arguments

- `connect --url <NOTIFICATION_URL> [--label <LABEL>]`<br>
  Reads the notification file, updates the repo with the latest snapshot, then the latest delta,
  and creates a new source record.
- `update  --source <SOURCE> [--label <LABEL>]`
  Reads the notification file, then updates the repo the latest delta,
- `list`
  Lists all sources in the repo.
- `rename --source <SOURCE> --label <FROM_LABEL> --to <TO_LABEL>`
  Replaces a label

A note about labels:
A label can be given to a source in order to track multiple sessions of the same IRR source.
This is useful when your repo can no longer be synchronized with a server, for example when the session
ID changes, or when your repo wasn't refreshed on time and loses sync.
In these cases you can use a label for each of the sessions to preserve the history, should you
wish to keep it. If you only want the latest version of each IRR source then you don't need labels.

GROW:

> a mirror server SHOULD remove all Delta Files older than 24 hours

# Quick set up

## PostgreSQL Database

### (Optional) Use Docker to run a local instance of Postgres

Set an environment variable `POSTGRES_HOST_AUTH_METHOD=trust` so you won't need to use passwords.
It's very insecure but also very handy for a local setup just to get things up and running. You
can always `pg_dump` the data and put it on a more securely configured server later.

    docker pull postgres:16
    docker run -d -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust --name db postgres:16

Or with podman

    podman pull docker.io/library/postgres:16
    podman run -d -p 5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust --name db postgres:16

From then on:

    docker stop db
    docker start db

### Create role and DB

Assuming your database is running on localhost...

    createuser -h localhost nrtm4
    createdb -h localhost -O nrtm4 nrtm4

    createuser -h localhost nrtm4_test
    createdb -h localhost -O nrtm4_test nrtm4_test

## Build

You'll need these tools:

- make, or `task`, which I'm moving towards.
- go 1.23+
- node 21+
- [tern](https://github.com/JackC/tern) v2.3.0 for PostgreSQL migrations

### Initialize schema

Edit `tern.conf` to contain the variables for your database, then...

    make migrate  # creates database schema and migrates it to the latest version

### Build targets

    make clean buildgo # creates a binary at ./cmd/nrtmclient/nrtmclient
    make clean testgo # uses a db to when testing. See above for PostgreSQL setup

The `run.sh` command should now be usable. See Usage above.

For development:

[This script](./scripts/pgdumpdata.sh) uses `pg_dump` to do a data-only dump of the
database. It excludes the schema version from the dump and encodes it in the output
file name. When restoring dumps, ensure the target schema matches the data dump
version.

Example usage:

    ./scripts/pgdumpdata.sh "-h localhost -U nrtm4" nrtm4dbhost

The result is a gzipped dump file which can be restored by piping the output to
a `psql` command in the usual way.

Other targets

    make emptydb  # wipes the table schema, including any data, ofc
    make rewinddb  # schema is set back one version
    make emptydb migratetest  # resets the test db
    make list # fill your boots

This builds the frontend as well, though I wouldn't bother until it can do cool stuff.

    make clean test
    make clean build

# Tips

Profile the code
https://granulate.io/blog/golang-profiling-basics-quick-tutorial/
