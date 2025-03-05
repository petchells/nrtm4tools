CREATE SEQUENCE _history_seq start 1
;

CREATE TABLE nrtm_rpslobject_history (
	id BIGINT NOT NULL PRIMARY KEY,
	seq BIGINT NOT NULL,
	stamp TIMESTAMP WITHOUT TIME ZONE,
	old_id BIGINT NOT NULL,
	object_type CHARACTER VARYING(255) NOT NULL,
	primary_key CHARACTER VARYING(255) NOT NULL,
	nrtm_source_id BIGINT NOT NULL,
	version INTEGER NOT NULL,
	rpsl TEXT NOT NULL
)
;

ALTER TABLE nrtm_rpslobject
DROP CONSTRAINT rpslobject__source__type__primary_key__to_version__uid
;
ALTER TABLE nrtm_rpslobject
DROP COLUMN to_version
;

ALTER TABLE nrtm_rpslobject
RENAME COLUMN from_version TO version
;

CREATE UNIQUE index nrtm_rpslobject_history__seq__idx ON nrtm_rpslobject_history(seq)
;

CREATE INDEX nrtm_rpslobject_history__source__idx ON nrtm_rpslobject_history(nrtm_source_id)
;

CREATE INDEX nrtm_rpslobject_history__type__key__idx ON nrtm_rpslobject_history(object_type, primary_key)
;

CREATE FUNCTION store_rpslobject_history () returns trigger AS $rpsl_history_recorder$
    DECLARE
        _seq bigint;
    BEGIN
        set timezone to 'UTC'; -- it should be anyway, but just in case
        SELECT nextval('_history_seq') INTO _seq;
        INSERT INTO nrtm_rpslobject_history
            (id, seq, stamp, old_id, object_type, primary_key, nrtm_source_id, version, rpsl)
        VALUES (
            id_generator(),
            _seq,
            now(),
            OLD.id,
            OLD.object_type,
            OLD.primary_key,
            OLD.nrtm_source_id,
            OLD.version,
            OLD.rpsl
        );
        RETURN NEW;
    END;
$rpsl_history_recorder$ language plpgsql
;

CREATE TRIGGER modify_rpsl_trigger before delete
OR
UPDATE ON nrtm_rpslobject FOR each ROW
EXECUTE function store_rpslobject_history ()
;

ALTER TABLE nrtm_rpslobject
DROP CONSTRAINT rpslobject__source__type__primary_key__from_version__uid;

---- create above / drop below ----

DROP TRIGGER modify_rpsl_trigger ON nrtm_rpslobject
;

DROP FUNCTION store_rpslobject_history
;

DROP INDEX nrtm_rpslobject_history__type__key__idx
;
DROP INDEX nrtm_rpslobject_history__source__idx
;
DROP index nrtm_rpslobject_history__seq__idx
;

ALTER TABLE nrtm_rpslobject
ADD COLUMN to_version INTEGER default 0 NOT NULL
;
alter table nrtm_rpslobject
add constraint rpslobject__source__type__primary_key__to_version__uid
    unique (nrtm_source_id, object_type, primary_key, to_version);

ALTER TABLE nrtm_rpslobject
RENAME COLUMN version TO from_version
;
ALTER TABLE nrtm_rpslobject
ADD CONSTRAINT rpslobject__source__type__primary_key__from_version__uid UNIQUE
(nrtm_source_id, object_type, primary_key, from_version);

DROP TABLE nrtm_rpslobject_history
;

DROP SEQUENCE _history_seq
;
