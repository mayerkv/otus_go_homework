package storage

import "time"

type EventID string

func (id EventID) String() string {
	return string(id)
}

type UserID string

func (id UserID) String() string {
	return string(id)
}

type Event struct {
	ID           EventID
	Title        string
	StartAt      time.Time
	EndAt        time.Time
	Description  string
	OwnerID      UserID
	NotifyBefore time.Duration
}
