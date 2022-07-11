package memorystorage

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage_Save(t *testing.T) {
	items := make(map[storage.EventID]storage.Event)
	mu := &sync.RWMutex{}
	store := &Storage{
		mu:    mu,
		items: items,
	}
	aEvent := storage.Event{
		ID:          storage.EventID(uuid.NewString()),
		Title:       "test",
		StartAt:     time.Now(),
		EndAt:       time.Now().Add(30 * time.Second),
		Description: "",
		OwnerID:     storage.UserID(uuid.NewString()),
		NotifyAt:    time.Time{},
	}

	err := store.Save(context.Background(), &aEvent)
	require.NoError(t, err)
	event, ok := items[aEvent.ID]
	require.True(t, ok)
	require.Equal(t, aEvent, event)
}

func TestStorage_FindByID(t *testing.T) {
	t.Run(
		"when event is exists, returns event", func(t *testing.T) {
			aEvent := storage.Event{
				ID:          storage.EventID(uuid.NewString()),
				Title:       "test",
				StartAt:     time.Now(),
				EndAt:       time.Now().Add(30 * time.Second),
				Description: "",
				OwnerID:     storage.UserID(uuid.NewString()),
				NotifyAt:    time.Time{},
			}
			store := &Storage{
				mu:    &sync.RWMutex{},
				items: map[storage.EventID]storage.Event{aEvent.ID: aEvent},
			}

			event, err := store.FindByID(context.Background(), aEvent.ID)
			require.NoError(t, err)
			require.Equal(t, aEvent, *event)
		},
	)

	t.Run(
		"when event does not exists, returns nil", func(t *testing.T) {
			aEvent := storage.Event{
				ID:          storage.EventID(uuid.NewString()),
				Title:       "test",
				StartAt:     time.Now(),
				EndAt:       time.Now().Add(30 * time.Second),
				Description: "",
				OwnerID:     storage.UserID(uuid.NewString()),
				NotifyAt:    time.Time{},
			}
			store := &Storage{
				mu:    &sync.RWMutex{},
				items: make(map[storage.EventID]storage.Event),
			}

			event, err := store.FindByID(context.Background(), aEvent.ID)
			require.NoError(t, err)
			require.Nil(t, event)
		},
	)

	t.Run(
		"concurrent access", func(t *testing.T) {
			aEvent := storage.Event{
				ID:          storage.EventID(uuid.NewString()),
				Title:       "test",
				StartAt:     time.Now(),
				EndAt:       time.Now().Add(30 * time.Second),
				Description: "",
				OwnerID:     storage.UserID(uuid.NewString()),
				NotifyAt:    time.Time{},
			}
			store := &Storage{
				mu:    &sync.RWMutex{},
				items: map[storage.EventID]storage.Event{aEvent.ID: aEvent},
			}

			wg := &sync.WaitGroup{}
			wg.Add(10)
			for i := 0; i < 10; i++ {
				go func(wg *sync.WaitGroup, aEvent storage.Event) {
					defer wg.Done()

					event, err := store.FindByID(context.Background(), aEvent.ID)
					require.NoError(t, err)
					require.Equal(t, aEvent, *event)
				}(wg, aEvent)
			}

			wg.Wait()
		},
	)
}

func TestStorage_Delete(t *testing.T) {
	t.Run(
		"when event is exists, then delete", func(t *testing.T) {
			aEvent := storage.Event{
				ID:          storage.EventID(uuid.NewString()),
				Title:       "test",
				StartAt:     time.Now(),
				EndAt:       time.Now().Add(30 * time.Second),
				Description: "",
				OwnerID:     storage.UserID(uuid.NewString()),
				NotifyAt:    time.Time{},
			}
			store := &Storage{
				mu:    &sync.RWMutex{},
				items: map[storage.EventID]storage.Event{aEvent.ID: aEvent},
			}
			err := store.Delete(context.Background(), &aEvent)
			require.NoError(t, err)

			_, ok := store.items[aEvent.ID]
			require.False(t, ok)
		},
	)

	t.Run(
		"when even does not exists, returns nil", func(t *testing.T) {
			aEvent := storage.Event{
				ID:          storage.EventID(uuid.NewString()),
				Title:       "test",
				StartAt:     time.Now(),
				EndAt:       time.Now().Add(30 * time.Second),
				Description: "",
				OwnerID:     storage.UserID(uuid.NewString()),
				NotifyAt:    time.Time{},
			}
			store := &Storage{
				mu:    &sync.RWMutex{},
				items: make(map[storage.EventID]storage.Event),
			}
			err := store.Delete(context.Background(), &aEvent)
			require.NoError(t, err)
			require.Len(t, store.items, 0)
		},
	)
}

func TestStorage_FindAllByUserIDAndPeriod(t *testing.T) {
	startAt, _ := time.Parse(time.RFC3339, "2006-01-01T10:00:00Z")
	endAt := startAt.Add(30 * time.Minute)
	aEvent := storage.Event{
		ID:          storage.EventID(uuid.NewString()),
		Title:       "test",
		StartAt:     startAt,
		EndAt:       endAt,
		Description: "",
		OwnerID:     storage.UserID(uuid.NewString()),
		NotifyAt:    time.Time{},
	}

	items := map[storage.EventID]storage.Event{
		aEvent.ID: aEvent,
	}
	mu := &sync.RWMutex{}
	store := &Storage{
		mu:    mu,
		items: items,
	}

	t.Run(
		"in range and existent user id", func(t *testing.T) {
			events, err := store.FindAllByUserIDAndPeriod(
				context.Background(), aEvent.OwnerID, startAt.Add(10*time.Minute), startAt.Add(24*time.Hour),
			)
			require.NoError(t, err)
			require.Equal(t, []storage.Event{aEvent}, events)
		},
	)

	t.Run(
		"in range and nonexistent user", func(t *testing.T) {
			events, err := store.FindAllByUserIDAndPeriod(
				context.Background(), "nonexistent_user_id", startAt.Add(10*time.Minute), startAt.Add(24*time.Hour),
			)
			require.NoError(t, err)
			require.Empty(t, events)
		},
	)

	t.Run(
		"not in range and existent user", func(t *testing.T) {
			events, err := store.FindAllByUserIDAndPeriod(
				context.Background(), aEvent.OwnerID, startAt.Add(30*time.Minute+time.Millisecond),
				startAt.Add(24*time.Hour),
			)
			require.NoError(t, err)
			require.Empty(t, events)
		},
	)
}

func TestStorage_HasByUserIDAndPeriod(t *testing.T) {
	startAt, _ := time.Parse(time.RFC3339, "2006-01-01T10:00:00Z")
	endAt := startAt.Add(30 * time.Minute)
	aEvent := storage.Event{
		ID:          storage.EventID(uuid.NewString()),
		Title:       "test",
		StartAt:     startAt,
		EndAt:       endAt,
		Description: "",
		OwnerID:     storage.UserID(uuid.NewString()),
		NotifyAt:    time.Time{},
	}

	items := map[storage.EventID]storage.Event{
		aEvent.ID: aEvent,
	}
	mu := &sync.RWMutex{}
	store := &Storage{
		mu:    mu,
		items: items,
	}

	t.Run(
		"in range and existent user id", func(t *testing.T) {
			ok, err := store.HasByUserIDAndPeriod(
				context.Background(),
				aEvent.OwnerID,
				startAt.Add(10*time.Minute),
				startAt.Add(24*time.Hour),
			)
			require.NoError(t, err)
			require.True(t, ok)
		},
	)

	t.Run(
		"in range and nonexistent user", func(t *testing.T) {
			ok, err := store.HasByUserIDAndPeriod(
				context.Background(),
				"nonexistent_user_id",
				startAt.Add(10*time.Minute),
				startAt.Add(24*time.Hour),
			)
			require.NoError(t, err)
			require.False(t, ok)
		},
	)

	t.Run(
		"not in range and existent user", func(t *testing.T) {
			ok, err := store.HasByUserIDAndPeriod(
				context.Background(),
				aEvent.OwnerID,
				startAt.Add(30*time.Minute+time.Millisecond),
				startAt.Add(24*time.Hour),
			)
			require.NoError(t, err)
			require.False(t, ok)
		},
	)
}

func TestStorage_NextID(t *testing.T) {
	store := &Storage{
		mu:    &sync.RWMutex{},
		items: map[storage.EventID]storage.Event{},
	}
	id, err := store.NextID(context.Background())
	require.NoError(t, err)

	u, err := uuid.Parse(id.String())
	require.NoError(t, err)
	require.Equal(t, id.String(), u.String())
}

func TestNew(t *testing.T) {
	expected := &Storage{&sync.RWMutex{}, map[storage.EventID]storage.Event{}}
	require.Equal(t, expected, New())
}
