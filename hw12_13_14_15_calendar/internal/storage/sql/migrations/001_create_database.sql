--+goose Up
create user igor with encrypted password 'igor';
create database if not exists calendardb;
-- +goose Down
drop database if exists calendardb;
drop user if exists igor;