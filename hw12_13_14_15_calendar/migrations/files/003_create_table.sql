-- +goose Up
create table if not exists events (
id serial primary key,
title text,
start timestamp with time zone not null,
duration bigint,
descr text,
notification bigint,
scheduled boolean
);

grant all privileges on table events to igor;
grant usage, select on all sequences in schema public to igor;
grant usage, update on all sequences in schema public to igor;

--+goose Down
-- drop table if exists events;