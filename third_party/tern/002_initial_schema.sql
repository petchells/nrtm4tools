create table nrtmstate (
	id bigint not null,
	source varchar(255) not null,
	version integer not null,
	url text not null,
	type varchar(255) not null,
	file_name text not null,
	created timestamp without time zone not null,
	
	constraint nrtmstate_pk primary key (id)
);

create table rpslobject (
	id bigint not null,
	object_type varchar(255) not null,
	rpsl text not null,
	source varchar(255) not null,
	primary_key varchar(255) not null,
	nrtmstate_id bigint not null,
	created timestamp without time zone not null,
	updated timestamp without time zone not null,

	constraint rpslobject_pk primary key (id),
	constraint rpslobject__source__primary_key_uid unique (source, object_type, primary_key),
	constraint rpslobject__nrtmstate_fk foreign key (nrtmstate_id) references nrtmstate(id)
);

create index nrtmstate__source_version_idx on nrtmstate(source, version);

---- create above / drop below ----
drop index nrtmstate__source_version_idx;
drop table rpslobject;
drop table nrtmstate;
