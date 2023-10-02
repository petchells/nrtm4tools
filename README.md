# NRTMv4 Client

## Introduction

nrtm4client is a tool for communicating with an [NRTMv4 server](https://www.ietf.org/archive/id/draft-ietf-grow-nrtm-v4-02.html). It can synchronize NTRM sources with a database (PostgreSQL currently) and validate the synchronization status of one or more NRTM servers.

## Build

* make
* go 1.20+
* [tern](https://github.com/JackC/tern) v1.13.0 for PostgreSQL migrations


## Database

### Set up a PostgreSQL database

	createuser -h localhost nrtm4
	createdb -h localhost -O nrtm4 nrtm4
	createuser -h localhost nrtm4_test
	createdb -h localhost -O nrtm4_test nrtm_test

link to docker

* Create role and DB
* Configure and run Tern to migrate db
* Script for doing this on a non-dev machine
* Configure Env/Flags PG_DATABASE_URL

### Build and dev targets

## Running nrtm4client

Create a directory, e.g. `$HOME/tmp/RIPE` to store downloaded files, then run it with these
environment variables (assumes the nrtm4client binary is at $HOME/Projects/nrtm4client/cmd/nrtm4client/nrtm4client):

	PG_DATABASE_URL="postgres://nrtm4:nrtm4@localhost:5432/nrtm4?sslmode=disable" \
	NRTM4_NOTIFICATION_URL="https://nrtm-rc.db.ripe.net/nrtmv4/RIPE/update-notification-file.json" \
	NRTM4_FILE_PATH=$HOME/tmp/RIPE \
	BOLT_DATABASE_PATH=$HOME/tmp/nrtm4.bbolt.db \
	$HOME/Projects/nrtm4client/cmd/nrtm4client/nrtm4client

Describe modes: Syncing and validating

* Environment variables
* Command line flags
