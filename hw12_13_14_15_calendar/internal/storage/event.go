package storage

import "time"

//easyjson:json
type Event struct {
	ID               int
	Title            string
	Start            time.Time
	Duration         int // think 'bout converting to time.Duration
	Description      string
	NotificationTime int // same as Duration time.Duration
	Scheduled        bool
}
