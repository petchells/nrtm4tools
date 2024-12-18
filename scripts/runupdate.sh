#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

varsfile="${SCRIPT_DIR}/env.dev.conf"

if [ ! -f $varsfile ];
then
	echo "Missing file: ${varsfile}"
	echo
	echo "Copy the file env.dev.example.conf to this location, then edit the values"
	echo "to match your system."
	exit 1
fi

env $(cat "$varsfile" | xargs) go run \
	cmd/nrtm4client/main.go update \
	--source ripe

