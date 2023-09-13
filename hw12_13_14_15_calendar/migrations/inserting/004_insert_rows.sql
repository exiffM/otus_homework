-- +goose Up
insert into events (title, start, duration, descr, notification, scheduled)
values ('Meeting', '2023-09-13 13:30', 30, 'Meeting with Fred Durst', 15, false);

insert into events (title, start, duration, descr, notification, scheduled)
values ('Meeting', '2023-09-13 13:31', 30, 'Meeting with Till Lindemann', 15, false);

insert into events (title, start, duration, descr, notification, scheduled)
values ('Meeting', '2023-09-13 13:32', 30, 'Meeting with Valentin Dyadka', 15, false);
-- +goose Down
drop table if exists events;