drop database if exists calendardb;

create database calendardb;

create user igor with encrypted password 'igor';

grant all privileges on database calendardb to igor;

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