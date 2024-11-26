# NRTMv4 Client

## Introduction

nrtm4client is a tool for communicating with an [NRTMv4 server](https://github.com/mxsasha/nrtmv4). It can synchronize NTRM sources with a database (PostgreSQL currently) and validate the synchronization status of one or more NRTM servers.

## Build

* make
* go 1.20+
* [tern](https://github.com/JackC/tern) v1.13.0 for PostgreSQL migrations


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

Create a directory, e.g. `$HOME/tmp/RIPE` to store downloaded files, then run it with these
environment variables (assumes the nrtm4client binary is at $HOME/Projects/nrtm4client/cmd/nrtm4client/nrtm4client):

	PG_DATABASE_URL="postgres://nrtm4:nrtm4@localhost:5432/nrtm4?sslmode=disable" \
	NRTM4_FILE_PATH=$HOME/tmp/RIPE \
	BOLT_DATABASE_PATH=$HOME/tmp/nrtm4.bbolt.db \
	$HOME/Projects/nrtm4client/cmd/nrtm4client/nrtm4client

Describe modes: Syncing and validating

* Environment variables
* Command line flags

# Tips

Profile the code
https://granulate.io/blog/golang-profiling-basics-quick-tutorial/
