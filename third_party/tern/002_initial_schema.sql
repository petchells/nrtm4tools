create table rpslobject (
	id bigint not null,
	created timestamp without time zone not null,
	updated timestamp without time zone not null,
	version bigint not null,
	rpsl text not null,
	source varchar(255) not null,
	primary_key varchar(255) not null,

	constraint rpslobject_pk primary key (id),
	constraint rpslobject__source__primary_key_uid unique (source, primary_key),
	unique (source, primary_key)
);

create table nrtmstate (
	id bigint not null,
	created timestamp without time zone not null,
	url text not null,
	is_delta boolean not null,
	delta text not null,
	source varchar(255) not null,
	
	constraint nrtmstate_pk primary key (id)
);
---- create above / drop below ----
drop table nrtmstate;
drop table rpslobject;