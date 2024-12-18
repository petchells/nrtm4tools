# NRTMv4 Client

## Introduction

nrtm4client is a tool for communicating with an [NRTMv4 server](https://github.com/mxsasha/nrtmv4).
It retrieves IRR data from NTRM mirror servers and stores them in a database (only PostgreSQL
currently). History is maintained and can be queried

## Development Status

The `main` branch supports `connect` and `update` commands. See [run.example.sh](./scripts/run.example.sh)
for example uses.

To try it out yourself follow the build steps below, then set up a PostgreSQL data before
making a connection to an NRTMv4 mirror server.

## Usage

The `connect` command downloads a snapshot file, inserts RPSL objects into the database, then
applies successive deltas until the database is up to date.

After that use the `update` command to download the
delta files and re-synchronize the database with the NRTM mirror server. If you want to
keep the history of changes to IRR records over time, you'll need to update regularly -- mirror
servers remove old delta files so you should get them while they're available.

## Build

- make
- go 1.23+
- [tern](https://github.com/JackC/tern) v2.3.0 for PostgreSQL migrations

## PostgreSQL Database

### Create role and DB

Assuming your database is running on localhost...

    createuser -h localhost nrtm4
    createdb -h localhost -O nrtm4 nrtm4
    createuser -h localhost nrtm4_test
    createdb -h localhost -O nrtm4_test nrtm4_test

### Initialize schema

    make migrate

## Running nrtm4client

Create a directory, e.g. `$HOME/nrtm4/RIPE` to store downloaded files, then run it with these
environment variables (assumes the nrtm4client binary is at $HOME/Projects/nrtm4client/cmd/nrtm4client/nrtm4client):

    envvars="PG_DATABASE_URL=postgres://nrtm4:nrtm4@localhost:5432/nrtm4?sslmode=disable \
        NRTM4_FILE_PATH=$HOME/tmp/RIPE"
    env ${envvars} ./cmd/nrtm4client/nrtm4client

Describe modes: Syncing and validating

- Environment variables
- Command line flags

# Tips

Profile the code
https://granulate.io/blog/golang-profiling-basics-quick-tutorial/
