#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

if [ -z "$SCRIPT_DIR" ];then
	echo "Cannot determine directory"
	exit 3
fi

if [ ! -d "$SCRIPT_DIR" ];then
	echo "Directory does not exist, or is not a directory: $SCRIPT_DIR"
	exit 2
fi

varsfile="${SCRIPT_DIR}/env.dev.conf"

if [ ! -f $varsfile ];
then
	echo "Missing file: ${varsfile}"
	echo
	echo "Copy the file env.example.conf to ${varsfile}, then edit the values"
	echo "to suit your system."
	exit 1
fi

cd "${SCRIPT_DIR}"/.. && \
env $(cat "$varsfile" | xargs) go run \
	cmd/nrtm4client/main.go "$@"
