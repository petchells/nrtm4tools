CREATE TABLE nrtm_source (
	id BIGINT NOT NULL,
	source VARCHAR(255) NOT NULL,
	session_id VARCHAR(255) NOT NULL,
	VERSION INTEGER NOT NULL,
	notification_url TEXT NOT NULL,
	label VARCHAR(255) NOT NULL,
	status VARCHAR(255) NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	CONSTRAINT nrtm_source__pk PRIMARY KEY (id),
	CONSTRAINT nrtm_source__source__label__uid UNIQUE (notification_url, label)
);

CREATE TABLE nrtm_notification (
	id BIGINT NOT NULL,
	VERSION INTEGER NOT NULL,
	source_id BIGINT NOT NULL,
	payload jsonb NOT NULL,
	created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	CONSTRAINT nrtm_notification__pk PRIMARY KEY (id),
	CONSTRAINT nrtm_notification__nrtm_source__fk FOREIGN key (source_id) REFERENCES nrtm_source (id)
);

CREATE INDEX nrtm_notification__version__idx ON nrtm_notification (source_id, VERSION);

CREATE TABLE nrtm_rpslobject (
	id BIGINT NOT NULL,
	object_type VARCHAR(255) NOT NULL,
	primary_key VARCHAR(255) NOT NULL,
	source_id BIGINT NOT NULL,
	VERSION INTEGER NOT NULL,
	rpsl TEXT NOT NULL,
	CONSTRAINT rpslobject__pk PRIMARY KEY (id),
	CONSTRAINT rpslobject__nrtm_source__fk FOREIGN key (source_id) REFERENCES nrtm_source (id),
	CONSTRAINT rpslobject__source__type__primary_key__uid UNIQUE (source_id, object_type, primary_key)
);

CREATE INDEX rpslobject__type__primary_key__idx ON nrtm_rpslobject (object_type, primary_key);

CREATE TABLE nrtm_rpslobject_history (
	id BIGINT NOT NULL PRIMARY KEY,
	seq BIGINT NOT NULL,
	stamp TIMESTAMP WITHOUT TIME ZONE,
	original_id BIGINT NOT NULL,
	object_type CHARACTER VARYING(255) NOT NULL,
	primary_key CHARACTER VARYING(255) NOT NULL,
	source_id BIGINT NOT NULL,
	VERSION INTEGER NOT NULL,
	rpsl TEXT NOT NULL
);

CREATE UNIQUE index nrtm_rpslobject_history__seq__idx ON nrtm_rpslobject_history (seq);

CREATE INDEX nrtm_rpslobject_history__source__idx ON nrtm_rpslobject_history (source_id);

CREATE INDEX nrtm_rpslobject_history__type__key__idx ON nrtm_rpslobject_history (object_type, primary_key);

CREATE FUNCTION store_rpslobject_history () returns trigger AS $rpsl_history_recorder$
    DECLARE
        _seq bigint;
    BEGIN
        set timezone to 'UTC'; -- it should be anyway, but just in case
        SELECT nextval('_history_seq') INTO _seq;
        INSERT INTO nrtm_rpslobject_history
            (id, seq, stamp, original_id, object_type, primary_key, source_id, version, rpsl)
        VALUES (
            id_generator(),
            _seq,
            now(),
            OLD.id,
            OLD.object_type,
            OLD.primary_key,
            OLD.source_id,
            OLD.version,
            OLD.rpsl
        );
        RETURN NEW;
    END;
$rpsl_history_recorder$ language plpgsql;

CREATE TRIGGER modify_rpsl_trigger before delete
OR
UPDATE ON nrtm_rpslobject FOR each ROW
EXECUTE function store_rpslobject_history ();

-----------------------------------
---- create above / drop below ----
-----------------------------------
DROP TRIGGER modify_rpsl_trigger ON nrtm_rpslobject;

DROP FUNCTION store_rpslobject_history;

DROP TABLE nrtm_rpslobject_history;

DROP TABLE nrtm_rpslobject;

DROP TABLE nrtm_notification;

DROP TABLE nrtm_source;
