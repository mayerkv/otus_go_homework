package app

import (
	"context"
	"errors"
	"time"

	"github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrDateBusy       = errors.New("date is busy")
	ErrEventNotExists = errors.New("event not exists")
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
	FindAllByUserIDAndPeriod(
		ctx context.Context, ownerID storage.UserID,
		from time.Time, to time.Time,
	) ([]storage.Event, error)
	HasByUserIDAndPeriod(ctx context.Context, ownerID storage.UserID, from time.Time, to time.Time) (bool, error)
	HasByUserIDAndPeriodForUpdate(
		ctx context.Context,
		forUpdate storage.EventID,
		ownerID storage.UserID,
		from time.Time,
		to time.Time,
	) (bool, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(
	ctx context.Context,
	title, description, ownerID string,
	startAt, endAt time.Time,
	notifyThreshold time.Duration,
) (storage.EventID, error) {
	isBusy, err := a.storage.HasByUserIDAndPeriod(ctx, storage.UserID(ownerID), startAt, endAt)
	if err != nil {
		return "", err
	}

	if isBusy {
		return "", ErrDateBusy
	}

	id, err := a.storage.NextID(ctx)
	if err != nil {
		return "", err
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

	if err := a.storage.Save(ctx, event); err != nil {
		return "", err
	}

	return event.ID, nil
}

func (a *App) UpdateEvent(
	ctx context.Context,
	eventID, title, description, ownerID string,
	startAt, endAt time.Time,
	notifyThreshold time.Duration,
) error {
	event, err := a.storage.FindByID(ctx, storage.EventID(eventID))
	if err != nil {
		return err
	}
	if event == nil {
		return ErrEventNotExists
	}

	event.Title = title
	event.Description = description
	event.OwnerID = storage.UserID(ownerID)
	event.StartAt = startAt
	event.EndAt = endAt
	event.NotifyAt = startAt.Truncate(notifyThreshold)

	isBusy, err := a.storage.HasByUserIDAndPeriodForUpdate(
		ctx,
		event.ID,
		event.OwnerID,
		event.StartAt,
		event.EndAt,
	)
	if err != nil {
		return err
	}
	if isBusy {
		return ErrDateBusy
	}

	return a.storage.Save(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	event, err := a.storage.FindByID(ctx, storage.EventID(id))
	if err != nil {
		return err
	}
	if event == nil {
		return ErrEventNotExists
	}
	return a.storage.Delete(ctx, event)
}

func (a *App) GetEventList(ctx context.Context, ownerID string, from, to time.Time) ([]storage.Event, error) {
	return a.storage.FindAllByUserIDAndPeriod(ctx, storage.UserID(ownerID), from, to)
}
