create sequence _history_seq start 1;

create table history_nrtm_rpslobject (
    id bigint not null primary_key,
    seq bigint not null,
    `timestamp` timestamp without timezone,
	old_id bigint not null,
    object_type character varying(255) NOT NULL,
    primary_key character varying(255) NOT NULL,
    nrtm_source_id bigint NOT NULL,
    from_version integer NOT NULL,
    to_version integer NOT NULL,
    rpsl text NOT NULL
);

alter table nrtm_rpslobject drop column to_version;
alter table nrtm_rpslobject column to_version version;

create unique index history_nrtm_rpslobject_seq__idx on history_nrtm_rpslobject(seq);
create index history_nrtm_rpslobject_source__idx on history_nrtm_rpslobject(source);
create index history_nrtm_rpslobject_type_key__idx on history_nrtm_rpslobject(object_type, primary_key);

CREATE FUNCTION update_rpslobject_history() RETURNS trigger AS $rpsl_history_recorder$
    BEGIN
        set timezone to 'UTC'; -- it should be anyway, but just in case
        SELECT nextval('_history_seq') INTO seq_id;
        INSERT INTO history_nrtm_rpslobject
            (seq, id, seq, object_type, primary_key, nrtm_source_id, from_version, to_version, rpsl)
        VALUES (
            id_generator(),
            seq_id,
            now(),
            OLD.id,
            OLD.object_type,
            OLD.primary_key,
            OLD.nrtm_source_id,
            OLD.version,
            NEW.version,
            OLD.rpsl
        );
        RETURN NEW;
    END;
$rpsl_history_recorder$ LANGUAGE plpgsql;

CREATE TRIGGER update_rpsl_trigger BEFORE DELETE OR UPDATE ON nrtm_rpslobject
    FOR EACH ROW EXECUTE FUNCTION update_rpslobject_history();

---- create above / drop below ----

drop trigger update_rpsl_trigger;
drop function update_rpslobject_history;
drop index history_nrtm_rpslobject_type_key__idx;
drop index history_nrtm_rpslobject_source__idx;
drop index history_nrtm_rpslobject_seq__idx;
drop sequence _history_seq;