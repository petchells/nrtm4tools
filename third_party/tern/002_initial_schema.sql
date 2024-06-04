create table nrtm_source (
	id bigint not null,
	source varchar(255) not null,
	session_id varchar(255) not null,
	version integer not null,
	notification_url text not null,
	label varchar(255) not null,
	created timestamp without time zone not null,

	constraint nrtm_source__pk primary key (id),
	constraint nrtm_source__source__label__uid unique (notification_url, label)
);

create table nrtm_notification (
	id bigint not null,
	version integer not null,
	nrtm_source_id bigint not null,
	payload jsonb not null,
	created timestamp without time zone not null,

	constraint nrtm_notification__pk primary key (id),
	constraint nrtm_notification__nrtm_source__fk foreign key(nrtm_source_id) references nrtm_source(id)
);

create index nrtm_notification__version__idx on nrtm_notification(version);

create table nrtm_file (
	id bigint not null,
	version integer not null,
	type varchar(255) not null,
	url text not null,
	file_name text not null,
	nrtm_source_id bigint not null,
	created timestamp without time zone not null,

	constraint nrtm_file__pk primary key (id),
	constraint nrtm_file__nrtm_source__fk foreign key(nrtm_source_id) references nrtm_source(id)
);

create index nrtm_file__source_version_idx on nrtm_file(nrtm_source_id, version);

create table nrtm_rpslobject (
	id bigint not null,
	object_type varchar(255) not null,
	primary_key varchar(255) not null,
	nrtm_source_id bigint not null,
	from_version integer not null,
	to_version integer not null,
	rpsl text not null,

	constraint rpslobject__pk primary key (id),
	constraint rpslobject__nrtm_source__fk foreign key (nrtm_source_id) references nrtm_source(id),
	constraint rpslobject__nrtm_source__object_type__primary_key__from_version__uid unique (nrtm_source_id, object_type, primary_key, from_version),
	constraint rpslobject__nrtm_source__object_type__primary_key__to_version__uid unique (nrtm_source_id, object_type, primary_key, to_version)
);

create index rpslobject__primary_key__idx on nrtm_rpslobject(upper(primary_key));

---- create above / drop below ----

drop table nrtm_rpslobject;
drop table nrtm_file;
drop table nrtm_notification;
drop table nrtm_source;
