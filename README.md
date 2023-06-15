# NRTMv4 Client

## Introduction

nrtm4client is a tool for communicating with an [NRTMv4 server](https://www.ietf.org/archive/id/draft-ietf-grow-nrtm-v4-02.html). It can synchronize NTRM sources with a database (PostgreSQL currently) and validate the synchronization status of one or more NRTM servers.

## Build

* make
* go 1.20
* [tern](https://github.com/JackC/tern) v1.13.0 for PostgreSQL migrations

### Set up a PostgreSQL database

link to docker

### Database

PostgreSQL

* Create role and DB
* Configure and run Tern to migrate db

### Build and dev targets

## Running nrtm4client

Describe modes: Syncing and validating

* Environment variables
* Command line flags


