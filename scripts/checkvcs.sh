#!/bin/sh

UPSTREAM=${1:-'@{u}'}
LOCAL=$(git rev-parse @)
REMOTE=$(git rev-parse "$UPSTREAM")
BASE=$(git merge-base @ "$UPSTREAM")

if [ $LOCAL = $REMOTE ]; then
	if git diff-index --quiet HEAD --; then
		exit 0
	else
		exit 20
	fi
elif [ $LOCAL = $BASE ]; then
	exit 21
elif [ $REMOTE = $BASE ]; then
	exit 22
fi
exit 23
