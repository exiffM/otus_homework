create schema if not exists calendardb;

drop schema calendardb cascade;

drop table if exists events;

create table events (
id serial primary key,
title text,
start timestamp with time zone not null,
duration bigint,
descr text,
notification bigint
);

grant all privileges on table events to igor;
grant usage, select on all sequences in schema public to igor;
grant usage, update on all sequences in schema public to igor;