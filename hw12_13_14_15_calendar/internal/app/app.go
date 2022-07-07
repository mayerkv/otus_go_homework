package app

import (
	"context"
	"errors"
	"time"

	"github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrDateBusy = errors.New("date is busy")
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	NextID(ctx context.Context) (storage.EventID, error)
	Save(ctx context.Context, event *storage.Event) error
	FindByID(ctx context.Context, eventID storage.EventID) (*storage.Event, error)
	Delete(ctx context.Context, event *storage.Event) error
	FindAllByUserIDAndPeriod(ctx context.Context, userID storage.UserID,
		from time.Time, to time.Time) ([]storage.Event, error)
	HasByUserIDAndPeriod(ctx context.Context, userID storage.UserID, from time.Time, to time.Time) (bool, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(ctx context.Context, title, description, ownerID string, startAt, endAt time.Time, notifyThreshold time.Duration) error {
	isBusy, err := a.storage.HasByUserIDAndPeriod(ctx, storage.UserID(ownerID), startAt, endAt)
	if err != nil {
		return err
	}

	if isBusy {
		return ErrDateBusy
	}

	id, err := a.storage.NextID(ctx)
	if err != nil {
		return err
	}

	event := &storage.Event{
		ID:          id,
		Title:       title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: description,
		OwnerID:     storage.UserID(ownerID),
		NotifyAt:    startAt.Truncate(notifyThreshold),
	}

	return a.storage.Save(ctx, event)
}
