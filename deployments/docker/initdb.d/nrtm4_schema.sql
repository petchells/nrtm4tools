--
-- PostgreSQL database dump
--

-- Dumped from database version 14.12 (Debian 14.12-1.pgdg120+1)
-- Dumped by pg_dump version 17.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO postgres;

--
-- Name: id_generator(); Type: FUNCTION; Schema: public; Owner: nrtm4
--

CREATE FUNCTION public.id_generator(OUT result bigint) RETURNS bigint
    LANGUAGE plpgsql
    AS $$
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
$$;


ALTER FUNCTION public.id_generator(OUT result bigint) OWNER TO nrtm4;

--
-- Name: store_rpslobject_history(); Type: FUNCTION; Schema: public; Owner: nrtm4
--

CREATE FUNCTION public.store_rpslobject_history() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
    DECLARE
        _seq bigint;
    BEGIN
        set timezone to 'UTC'; -- it should be anyway, but just in case
        SELECT nextval('_history_seq') INTO _seq;
        INSERT INTO nrtm_rpslobject_history
            (id, seq, stamp, old_id, object_type, primary_key, source_id, version, rpsl)
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
$$;


ALTER FUNCTION public.store_rpslobject_history() OWNER TO nrtm4;

--
-- Name: _history_seq; Type: SEQUENCE; Schema: public; Owner: nrtm4
--

CREATE SEQUENCE public._history_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public._history_seq OWNER TO nrtm4;

--
-- Name: _seq; Type: SEQUENCE; Schema: public; Owner: nrtm4
--

CREATE SEQUENCE public._seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public._seq OWNER TO nrtm4;

SET default_table_access_method = heap;

--
-- Name: nrtm_notification; Type: TABLE; Schema: public; Owner: nrtm4
--

CREATE TABLE public.nrtm_notification (
    id bigint NOT NULL,
    version integer NOT NULL,
    source_id bigint NOT NULL,
    payload jsonb NOT NULL,
    created timestamp without time zone NOT NULL
);


ALTER TABLE public.nrtm_notification OWNER TO nrtm4;

--
-- Name: nrtm_rpslobject; Type: TABLE; Schema: public; Owner: nrtm4
--

CREATE TABLE public.nrtm_rpslobject (
    id bigint NOT NULL,
    object_type character varying(255) NOT NULL,
    primary_key character varying(255) NOT NULL,
    source_id bigint NOT NULL,
    version integer NOT NULL,
    rpsl text NOT NULL
);


ALTER TABLE public.nrtm_rpslobject OWNER TO nrtm4;

--
-- Name: nrtm_rpslobject_history; Type: TABLE; Schema: public; Owner: nrtm4
--

CREATE TABLE public.nrtm_rpslobject_history (
    id bigint NOT NULL,
    seq bigint NOT NULL,
    stamp timestamp without time zone,
    old_id bigint NOT NULL,
    object_type character varying(255) NOT NULL,
    primary_key character varying(255) NOT NULL,
    source_id bigint NOT NULL,
    version integer NOT NULL,
    rpsl text NOT NULL
);


ALTER TABLE public.nrtm_rpslobject_history OWNER TO nrtm4;

--
-- Name: nrtm_source; Type: TABLE; Schema: public; Owner: nrtm4
--

CREATE TABLE public.nrtm_source (
    id bigint NOT NULL,
    source character varying(255) NOT NULL,
    session_id character varying(255) NOT NULL,
    version integer NOT NULL,
    notification_url text NOT NULL,
    label character varying(255) NOT NULL,
    created timestamp without time zone NOT NULL
);


ALTER TABLE public.nrtm_source OWNER TO nrtm4;

--
-- Name: schema_version; Type: TABLE; Schema: public; Owner: nrtm4
--

CREATE TABLE public.schema_version (
    version integer NOT NULL
);


ALTER TABLE public.schema_version OWNER TO nrtm4;

--
-- Name: nrtm_notification nrtm_notification__pk; Type: CONSTRAINT; Schema: public; Owner: nrtm4
--

ALTER TABLE ONLY public.nrtm_notification
    ADD CONSTRAINT nrtm_notification__pk PRIMARY KEY (id);


--
-- Name: nrtm_rpslobject_history nrtm_rpslobject_history_pkey; Type: CONSTRAINT; Schema: public; Owner: nrtm4
--

ALTER TABLE ONLY public.nrtm_rpslobject_history
    ADD CONSTRAINT nrtm_rpslobject_history_pkey PRIMARY KEY (id);


--
-- Name: nrtm_source nrtm_source__pk; Type: CONSTRAINT; Schema: public; Owner: nrtm4
--

ALTER TABLE ONLY public.nrtm_source
    ADD CONSTRAINT nrtm_source__pk PRIMARY KEY (id);


--
-- Name: nrtm_source nrtm_source__source__label__uid; Type: CONSTRAINT; Schema: public; Owner: nrtm4
--

ALTER TABLE ONLY public.nrtm_source
    ADD CONSTRAINT nrtm_source__source__label__uid UNIQUE (notification_url, label);


--
-- Name: nrtm_rpslobject rpslobject__pk; Type: CONSTRAINT; Schema: public; Owner: nrtm4
--

ALTER TABLE ONLY public.nrtm_rpslobject
    ADD CONSTRAINT rpslobject__pk PRIMARY KEY (id);


--
-- Name: nrtm_rpslobject rpslobject__source__type__primary_key__version__uid; Type: CONSTRAINT; Schema: public; Owner: nrtm4
--

ALTER TABLE ONLY public.nrtm_rpslobject
    ADD CONSTRAINT rpslobject__source__type__primary_key__version__uid UNIQUE (source_id, object_type, primary_key, version);


--
-- Name: nrtm_notification__version__idx; Type: INDEX; Schema: public; Owner: nrtm4
--

CREATE INDEX nrtm_notification__version__idx ON public.nrtm_notification USING btree (source_id, version);


--
-- Name: nrtm_rpslobject_history__seq__idx; Type: INDEX; Schema: public; Owner: nrtm4
--

CREATE UNIQUE INDEX nrtm_rpslobject_history__seq__idx ON public.nrtm_rpslobject_history USING btree (seq);


--
-- Name: nrtm_rpslobject_history__source__idx; Type: INDEX; Schema: public; Owner: nrtm4
--

CREATE INDEX nrtm_rpslobject_history__source__idx ON public.nrtm_rpslobject_history USING btree (source_id);


--
-- Name: nrtm_rpslobject_history__type__key__idx; Type: INDEX; Schema: public; Owner: nrtm4
--

CREATE INDEX nrtm_rpslobject_history__type__key__idx ON public.nrtm_rpslobject_history USING btree (object_type, primary_key);


--
-- Name: rpslobject__type__primary_key__idx; Type: INDEX; Schema: public; Owner: nrtm4
--

CREATE INDEX rpslobject__type__primary_key__idx ON public.nrtm_rpslobject USING btree (object_type, primary_key);


--
-- Name: nrtm_rpslobject modify_rpsl_trigger; Type: TRIGGER; Schema: public; Owner: nrtm4
--

CREATE TRIGGER modify_rpsl_trigger BEFORE DELETE OR UPDATE ON public.nrtm_rpslobject FOR EACH ROW EXECUTE FUNCTION public.store_rpslobject_history();


--
-- Name: nrtm_notification nrtm_notification__nrtm_source__fk; Type: FK CONSTRAINT; Schema: public; Owner: nrtm4
--

ALTER TABLE ONLY public.nrtm_notification
    ADD CONSTRAINT nrtm_notification__nrtm_source__fk FOREIGN KEY (source_id) REFERENCES public.nrtm_source(id);


--
-- Name: nrtm_rpslobject rpslobject__nrtm_source__fk; Type: FK CONSTRAINT; Schema: public; Owner: nrtm4
--

ALTER TABLE ONLY public.nrtm_rpslobject
    ADD CONSTRAINT rpslobject__nrtm_source__fk FOREIGN KEY (source_id) REFERENCES public.nrtm_source(id);


--
-- PostgreSQL database dump complete
--

