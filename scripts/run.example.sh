#!/bin/sh

NRTM_FILE_DIR=/tmp/nrtm4

envvars="\
	PG_DATABASE_URL=postgres://nrtm4:nrtm4@localhost:5432/nrtm4?sslmode=disable \
	NRTM4_NOTIFICATION_URL='https://nrtm.db.ripe.net/nrtmv4/RIPE/update-notification-file.json \
	NRTM4_FILE_PATH=${NRTM_FILE_DIR} \
	BOLT_DATABASE_PATH=${NRTM_FILE_DIR} \
	"
env ${envvars} go run cmd/nrtm4client/main.go
