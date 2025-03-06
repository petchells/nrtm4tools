SET
	timezone TO 'UTC';

CREATE SEQUENCE _seq start 1;

CREATE SEQUENCE _history_seq start 1;

CREATE
OR REPLACE function id_generator (OUT result BIGINT) AS $$
DECLARE
    our_epoch bigint := 1741209445083;
    seq_id bigint;
    now_millis bigint;
    -- the id of this DB shard, must be set for each
    -- schema shard you have - you could pass this as a parameter too
    shard_id int := 1; -- up to 1024
BEGIN
    SELECT nextval('_seq') % 4096 INTO seq_id;
    SELECT FLOOR(EXTRACT(EPOCH FROM clock_timestamp()) * 1000) INTO now_millis;
    result := (now_millis - our_epoch) << 22;
    result := result | (shard_id << 12);
    result := result | (seq_id);
END;
$$ language plpgsql;

---- create above / drop below ----
DROP FUNCTION if EXISTS id_generator;

DROP SEQUENCE if EXISTS _history_seq;

DROP SEQUENCE if EXISTS _seq;