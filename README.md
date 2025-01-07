# NRTMv4 Client

## Introduction

nrtm4client is a tool for communicating with an [NRTMv4 server](https://github.com/mxsasha/nrtmv4).
It retrieves IRR data from NTRM mirror servers and stores them in a database (only PostgreSQL
currently). History is maintained and can be queried.

## Development Status

The `main` branch supports `connect` and `update` commands. See the `run*.sh` commands in
the `./scripts` directory. Create a file `./scripts/env.dev.conf`, set the vars to your
environment and you can use the scripts. Currently the only available source is RIPE, so
it's hard-coded.

## Set up environment variables

- NRTM4_FILE_PATH An empty directory where NRTMv4 snapshot and delta files will be stored.
- PG_DATABASE_URL Connection string to PostgreSQL database.

To try it out yourself follow the build steps below, then set up a PostgreSQL data before
making a connection to an NRTMv4 mirror server.

## Usage

The `connect` command downloads a snapshot file, inserts RPSL objects into the database, then
applies successive deltas until the database is up to date.

After that use the `update` command to download the
delta files and re-synchronize the database with the NRTM mirror server. If you want to
keep the history of changes to IRR records over time, you'll need to update regularly -- mirror
servers remove old delta files so you should get them while they're available.

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

### Initialize schema

    make migrate  # creates database schema. See `tern.conf`

## Build

- make
- go 1.23+
- [tern](https://github.com/JackC/tern) v2.3.0 for PostgreSQL migrations

  make clean testgo # uses a db to when testing. see below for PostgreSQL setup
  make clean buildgo # creates a binary at ./cmd/nrtmclient/nrtmclient

This builds the frontend as well, though I wouldn't bother until it can do cool stuff.

    make clean test
    make clean build

## Running nrtm4client

Create a directory, e.g. `$HOME/nrtm4/RIPE` to store downloaded files,
then copy file [./scripts/env.dev.example.conf] to ./scripts/env.dev.conf, and change the variables
to your system:

    PG_DATABASE_URL=postgres://nrtm4:nrtm4@localhost:5432/nrtm4?sslmode=disable
    NRTM4_FILE_PATH=$HOME/tmp/RIPE

Now run one of one of the `run*.sh` scripts in the [./scripts](./scripts) dir.

# Tips

Profile the code
https://granulate.io/blog/golang-profiling-basics-quick-tutorial/
