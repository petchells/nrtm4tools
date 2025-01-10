#!/bin/bash

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

echo $connect_opts
echo $machine_name

exit

schema_version=$(psql $connect_opts -P tuples_only -c 'select version from schema_version limit 1;' | grep -o '[0-9]\+')
TZ=Z
datestamp=$(date --rfc-3339='seconds'|grep -o '.\{19\}'|tr ' :' T-)
output_file_name="$machine_name"-data-filextract_v"$schema_version"_"$datestamp".dmp.sql.gz
pg_dump -t 'nrtm_*' -T 'schema_version' --data-only $connect_opts |gzip > $output_file_name
