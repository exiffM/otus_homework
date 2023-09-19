-- +goose Up
insert into events (title, start, duration, descr, notification, scheduled)
values ('Meeting', '2023-09-16 18:20', 30, 'Meeting with Fred Durst', 15, false);

insert into events (title, start, duration, descr, notification, scheduled)
values ('Meeting', '2023-09-16 18:21', 30, 'Meeting with Till Lindemann', 15, false);

insert into events (title, start, duration, descr, notification, scheduled)
values ('Meeting', '2023-09-16 18:22', 30, 'Meeting with Valentin Dyadka', 15, false);
-- +goose Down
-- drop table if exists events;