--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5 (Ubuntu 17.5-1.pgdg22.04+1)
-- Dumped by pg_dump version 17.5 (Ubuntu 17.5-1.pgdg22.04+1)

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
-- Name: timescaledb; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS timescaledb WITH SCHEMA public;


--
-- Name: EXTENSION timescaledb; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION timescaledb IS 'Enables scalable inserts and complex queries for time-series data (Community Edition)';


--
-- Name: timescaledb_toolkit; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS timescaledb_toolkit WITH SCHEMA public;


--
-- Name: EXTENSION timescaledb_toolkit; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION timescaledb_toolkit IS 'Library of analytical hyperfunctions, time-series pipelining, and other SQL utilities';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: ambient_station_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.ambient_station_data (
    id integer NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    date text,
    timezone text,
    date_utc integer,
    inside_temp_f double precision,
    inside_temp_f_feels_like double precision,
    inside_temp_c double precision,
    outside_temp_f double precision,
    outside_temp_f_feels_like double precision,
    inside_humidity double precision,
    outside_humidity double precision,
    inside_dew_point double precision,
    outside_dew_point double precision,
    barometric_relative double precision,
    barometric_absolute double precision,
    wind_direction double precision,
    wind_speed_mph double precision,
    wind_gust_mph double precision,
    max_daily_gust_mph double precision,
    hourly_rain_inches double precision,
    event_rain_inches double precision,
    daily_rain_inches double precision,
    weekly_rain_inches double precision,
    monthly_rain_inches double precision,
    total_rain_inches double precision,
    uv_index double precision,
    solar_radiation double precision,
    outside_batt_status integer,
    co2_bat_status integer,
    device_id integer,
    device_type_id integer
);


ALTER TABLE public.ambient_station_data OWNER TO postgres;

--
-- Name: ambient_station_data_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.ambient_station_data_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ambient_station_data_id_seq OWNER TO postgres;

--
-- Name: ambient_station_data_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.ambient_station_data_id_seq OWNED BY public.ambient_station_data.id;


--
-- Name: avtech_data; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.avtech_data (
    id integer NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    temp_f double precision NOT NULL,
    temp_c double precision NOT NULL,
    device_id integer NOT NULL,
    device_type_id integer NOT NULL
);


ALTER TABLE public.avtech_data OWNER TO postgres;

--
-- Name: avtech_data_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.avtech_data_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.avtech_data_id_seq OWNER TO postgres;

--
-- Name: avtech_data_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.avtech_data_id_seq OWNED BY public.avtech_data.id;


--
-- Name: collector_groups; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.collector_groups (
    id integer NOT NULL,
    group_name character varying NOT NULL,
    enabled boolean DEFAULT (1)::boolean NOT NULL,
    poll_interval_seconds integer DEFAULT 30
);


ALTER TABLE public.collector_groups OWNER TO postgres;

--
-- Name: COLUMN collector_groups.enabled; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collector_groups.enabled IS 'Global flag to enable collection for all devices in this group';


--
-- Name: COLUMN collector_groups.poll_interval_seconds; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.collector_groups.poll_interval_seconds IS 'The global poll interval for each device in the collector group';


--
-- Name: collector_groups_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.collector_groups_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.collector_groups_id_seq OWNER TO postgres;

--
-- Name: collector_groups_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.collector_groups_id_seq OWNED BY public.collector_groups.id;


--
-- Name: device_list; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.device_list (
    id integer NOT NULL,
    device_name character varying NOT NULL,
    location character varying DEFAULT ''::character varying,
    ip_address character varying DEFAULT ''::character varying,
    device_type_id integer NOT NULL,
    collector_group_id integer,
    enabled boolean DEFAULT (1)::boolean,
    poll_interval_seconds integer
);


ALTER TABLE public.device_list OWNER TO postgres;

--
-- Name: COLUMN device_list.enabled; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.device_list.enabled IS 'Enables data collection for this device';


--
-- Name: COLUMN device_list.poll_interval_seconds; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.device_list.poll_interval_seconds IS 'Poll interval in seconds for this device';


--
-- Name: device_list_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.device_list_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.device_list_id_seq OWNER TO postgres;

--
-- Name: device_list_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.device_list_id_seq OWNED BY public.device_list.id;


--
-- Name: device_types; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.device_types (
    id integer NOT NULL,
    device_type text NOT NULL
);


ALTER TABLE public.device_types OWNER TO postgres;

--
-- Name: ambient_station_data id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ambient_station_data ALTER COLUMN id SET DEFAULT nextval('public.ambient_station_data_id_seq'::regclass);


--
-- Name: avtech_data id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.avtech_data ALTER COLUMN id SET DEFAULT nextval('public.avtech_data_id_seq'::regclass);


--
-- Name: collector_groups id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collector_groups ALTER COLUMN id SET DEFAULT nextval('public.collector_groups_id_seq'::regclass);


--
-- Name: device_list id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.device_list ALTER COLUMN id SET DEFAULT nextval('public.device_list_id_seq'::regclass);


--
-- Name: ambient_station_data ambient_station_data_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.ambient_station_data
    ADD CONSTRAINT ambient_station_data_pk PRIMARY KEY (id);


--
-- Name: avtech_data avtech_data_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.avtech_data
    ADD CONSTRAINT avtech_data_pkey PRIMARY KEY (device_id, "timestamp");


--
-- Name: collector_groups collector_groups_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collector_groups
    ADD CONSTRAINT collector_groups_pk PRIMARY KEY (id);


--
-- Name: collector_groups collector_groups_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.collector_groups
    ADD CONSTRAINT collector_groups_unique UNIQUE (group_name);


--
-- Name: device_list device_list_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.device_list
    ADD CONSTRAINT device_list_pk PRIMARY KEY (id);


--
-- Name: device_list device_list_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.device_list
    ADD CONSTRAINT device_list_unique UNIQUE (device_name);


--
-- Name: device_types device_types_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.device_types
    ADD CONSTRAINT device_types_pk PRIMARY KEY (id);


--
-- Name: device_types device_types_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.device_types
    ADD CONSTRAINT device_types_unique UNIQUE (device_type);


--
-- Name: avtech_data_device_type_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX avtech_data_device_type_idx ON public.avtech_data USING btree (device_type_id);


--
-- Name: avtech_data_timestamp_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX avtech_data_timestamp_idx ON public.avtech_data USING btree ("timestamp" DESC);


--
-- Name: avtech_data ts_insert_blocker; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER ts_insert_blocker BEFORE INSERT ON public.avtech_data FOR EACH ROW EXECUTE FUNCTION _timescaledb_functions.insert_blocker();


--
-- Name: avtech_data avtech_data_device_types_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.avtech_data
    ADD CONSTRAINT avtech_data_device_types_fk FOREIGN KEY (device_type_id) REFERENCES public.device_types(id);


--
-- Name: device_list device_list_collector_groups_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.device_list
    ADD CONSTRAINT device_list_collector_groups_fk FOREIGN KEY (collector_group_id) REFERENCES public.collector_groups(id);


--
-- Name: device_list device_list_device_types_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.device_list
    ADD CONSTRAINT device_list_device_types_fk FOREIGN KEY (device_type_id) REFERENCES public.device_types(id);


--
-- PostgreSQL database dump complete
--

