package sqlstorage

import (
	"context"
	"database/sql"

	mdl "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq" // comment for justifying
)

type Storage struct {
	dsn    string
	dbConn *sql.DB
}

func New(dsn string) *Storage {
	return &Storage{dsn: dsn}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	_ = ctx
	db, err := sql.Open("postgres", s.dsn)
	if err != nil {
		return err
	}
	s.dbConn = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	_ = ctx
	if s.dbConn == nil {
		return nil
	}
	err := s.dbConn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) CreateEvent(event mdl.Event) error {
	sqlStatement := `
	INSERT INTO events(
		"title", "start", "duration", "descr", "notification"
	) values($1, $2, $3, $4, $5)
	RETURNING "id";`
	err := s.dbConn.QueryRow(sqlStatement,
		event.Title, event.Start, event.Duration, event.Description, event.NotificationTime,
	).Scan(&(event.ID))
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) SelectEvent(id int) (mdl.Event, error) {
	var event mdl.Event
	sqlStatement := `
	SELECT * FROM events WHERE "id"=$1`
	row := s.dbConn.QueryRow(sqlStatement, id)
	err := row.Scan(&(event.ID), &(event.Title), &(event.Start),
		&(event.Duration), &(event.Description), &(event.NotificationTime))
	if err != nil {
		return mdl.Event{}, err
	}
	return event, nil
}

func (s *Storage) UpdateEvent(event mdl.Event) error {
	sqlStatement := `
	UPDATE events 
	SET "title"=$1, "start"=$2, "duration"=$3, "descr"=$4, "notification"=$5
	WHERE "id"=$6;`
	_, err := s.dbConn.Exec(sqlStatement, event.Title, event.Start,
		event.Duration, event.Description, event.NotificationTime, event.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) DeleteEvent(id int) error {
	sqlStatement := `DELETE FROM events WHERE id=$1;`
	_, err := s.dbConn.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Events() ([]mdl.Event, error) {
	var event mdl.Event
	result := make([]mdl.Event, 0)
	sqlStatement := `
	SELECT * FROM events`
	rows, err := s.dbConn.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&event.ID, &event.Title, &event.Start, &event.Duration, &event.Description, &event.NotificationTime)
		if err != nil {
			return nil, err
		}
		result = append(result, event)
	}
	return result, nil
}
