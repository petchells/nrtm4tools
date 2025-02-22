CREATE SEQUENCE _history_seq start 1
;

CREATE TABLE history_nrtm_rpslobject (
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
DROP COLUMN to_version
;

ALTER TABLE nrtm_rpslobject
RENAME COLUMN from_version TO version
;

CREATE UNIQUE index history_nrtm_rpslobject_seq__idx ON history_nrtm_rpslobject (seq)
;

CREATE INDEX history_nrtm_rpslobject_source__idx ON history_nrtm_rpslobject (nrtm_source_id)
;

CREATE INDEX history_nrtm_rpslobject_type_key__idx ON history_nrtm_rpslobject (object_type, primary_key)
;

CREATE FUNCTION update_rpslobject_history () returns trigger AS $rpsl_history_recorder$
    DECLARE
        _seq bigint;
    BEGIN
        set timezone to 'UTC'; -- it should be anyway, but just in case
        SELECT nextval('_history_seq') INTO _seq;
        INSERT INTO history_nrtm_rpslobject
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

CREATE TRIGGER update_rpsl_trigger before delete
OR
UPDATE ON nrtm_rpslobject FOR each ROW
EXECUTE function update_rpslobject_history ()
;

---- create above / drop below ----

DROP TRIGGER update_rpsl_trigger ON nrtm_rpslobject
;

DROP FUNCTION update_rpslobject_history
;

DROP INDEX history_nrtm_rpslobject_type_key__idx
;

DROP INDEX history_nrtm_rpslobject_source__idx
;

DROP INDEX history_nrtm_rpslobject_seq__idx
;

ALTER TABLE nrtm_rpslobject
ADD COLUMN to_version INTEGER NOT NULL
;

ALTER TABLE nrtm_rpslobject
RENAME COLUMN version TO from_version
;

DROP TABLE history_nrtm_rpslobject
;

DROP SEQUENCE _history_seq
;
