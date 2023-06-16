# NRTMv4 Client

## Introduction

nrtm4client is a tool for communicating with an [NRTMv4 server](https://www.ietf.org/archive/id/draft-ietf-grow-nrtm-v4-02.html). It can synchronize NTRM sources with a database (PostgreSQL currently) and validate the synchronization status of one or more NRTM servers.

## Build

* make
* go 1.20+
* [tern](https://github.com/JackC/tern) v1.13.0 for PostgreSQL migrations


## Database

### Set up a PostgreSQL database

	createuser -h localhost nrtm4client
	createdb -h localhost -O nrtm4client nrtm4client
	createuser -h localhost nrtm4client_test
	createdb -h localhost -O nrtm4client_test nrtm4client_test

link to docker

* Create role and DB
* Configure and run Tern to migrate db
* Script for doing this on a non-dev machine
* Configure Env/Flags DATABASE_URL

### Build and dev targets

## Running nrtm4client

	"PG_DATEBASE_URL": "postgres://nrtm4client:nrtm4client@localhost:5432/nrtm4client?sslmode=disable",
	"NRTM4_NOTIFICATION_URL": "https://nrtm-rc.db.ripe.net/nrtmv4/RIPE/update-notification-file.json",
	"NRTM4_FILE_PATH": "/Users/etch/tmp/RIPE"

Describe modes: Syncing and validating

* Environment variables
* Command line flags
