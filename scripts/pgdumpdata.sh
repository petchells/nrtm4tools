#!/bin/bash


# pgdumpdata.sh - does a data-only dump of the tables in the targetted schema,
# excluding the table `schema_version`. Produces a gzipped dump file in the
# current working directory.
#
# E.g.
# connect_opts="-h localhost -U nrtm4"
# machine_name=bastion
connect_opts=
machine_name=

if [ -n "$1" ];then
    connect_opts=$1
else
    echo "Must provide connect options"
    exit 1
fi
if [ -n "$2" ];then
    machine_name=$2
else
    echo "Must provide machine name"
    exit 1
fi
if [ $# != 2 ];then
    echo "Invalid arguments"
    exit 1
fi

schema_version=$(psql $connect_opts -P tuples_only -c 'select version from schema_version limit 1;' | grep -o '[0-9]\+')
TZ=Z
datestamp=$(date --rfc-3339='seconds'|grep -o '.\{19\}'|tr ' :' T-)
output_file_name="$machine_name"-nrtm4data_v"$schema_version"_"$datestamp".dmp.sql.gz
pg_dump -t 'nrtm_*' -T 'schema_version' --data-only $connect_opts |gzip > $output_file_name
