package storage

import "time"

type EventID string

type UserID string

type Event struct {
	ID          EventID
	Title       string
	StartAt     time.Time
	EndAt       time.Time
	Description string
	OwnerID     UserID
	NotifyAt    time.Time
}
