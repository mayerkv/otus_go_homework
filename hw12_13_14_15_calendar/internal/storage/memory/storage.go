package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu    *sync.RWMutex
	items map[storage.EventID]storage.Event
}

func (s *Storage) NextID(ctx context.Context) (storage.EventID, error) {
	return storage.EventID(uuid.NewString()), nil
}

func (s *Storage) Save(ctx context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[event.ID] = *event
	return nil
}

func (s *Storage) FindByID(ctx context.Context, eventID storage.EventID) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if event, ok := s.items[eventID]; ok {
		return &event, nil
	}

	return nil, nil
}

func (s *Storage) Delete(ctx context.Context, event *storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, event.ID)
	return nil
}

func (s *Storage) FindAllByUserIDAndPeriod(
	ctx context.Context, ownerID storage.UserID,
	from, to time.Time,
) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events := make([]storage.Event, 0)
	for _, event := range s.items {
		if event.OwnerID != ownerID {
			continue
		}
		if s.inRange(event.StartAt, from, to) || s.inRange(event.EndAt, from, to) {
			events = append(events, event)
		}
	}
	return events, nil
}

func (s *Storage) HasByUserIDAndPeriod(ctx context.Context, ownerID storage.UserID, from, to time.Time) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, event := range s.items {
		if event.OwnerID != ownerID {
			continue
		}
		if s.inRange(event.StartAt, from, to) || s.inRange(event.EndAt, from, to) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Storage) HasByUserIDAndPeriodForUpdate(
	ctx context.Context,
	forUpdate storage.EventID,
	ownerID storage.UserID,
	from time.Time,
	to time.Time,
) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, event := range s.items {
		if event.OwnerID != ownerID {
			continue
		}
		if event.ID == forUpdate {
			if s.inRange(from, event.StartAt, event.EndAt) && s.inRange(to, event.StartAt, event.EndAt) {
				return false, nil
			}
		}
		if s.inRange(event.StartAt, from, to) || s.inRange(event.EndAt, from, to) {
			return true, nil
		}
	}
	return false, nil
}

func (s *Storage) inRange(date, from, to time.Time) bool {
	return date.Equal(from) || date.Equal(to) || (date.After(from) && date.Before(to))
}

func New() *Storage {
	return &Storage{
		mu:    &sync.RWMutex{},
		items: map[storage.EventID]storage.Event{},
	}
}
