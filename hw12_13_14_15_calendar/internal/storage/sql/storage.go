package sqlstorage

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib" // pg driver
	"github.com/mayerkv/otus_go_homework/hw12_13_14_15_calendar/internal/storage"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Storage struct {
	db              *sql.DB
	dsn             string
	maxOpenConns    int
	connMaxLifetime time.Duration
	maxIdleConns    int
	connMaxIdleTime time.Duration
}

func (s *Storage) NextID(ctx context.Context) (storage.EventID, error) {
	return storage.EventID(uuid.NewString()), nil
}

func (s *Storage) Save(ctx context.Context, event *storage.Event) error {
	_, err := s.db.ExecContext(ctx, saveQuery, event.ID, event.Title, event.StartAt, event.EndAt, event.Description, event.OwnerID, event.NotifyAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) FindByID(ctx context.Context, eventID storage.EventID) (*storage.Event, error) {
	row := s.db.QueryRowContext(ctx, selectQuery, eventID)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, nil
	}
	if row.Err() != nil {
		return nil, row.Err()
	}

	var event storage.Event
	if err := row.Scan(event.ID, event.Title, event.StartAt, event.EndAt, event.Description, event.NotifyAt); err != nil {
		return nil, err
	}

	return &event, nil
}

func (s *Storage) Delete(ctx context.Context, event *storage.Event) error {
	_, err := s.db.ExecContext(ctx, deleteQuery, event.ID)
	return err
}

func (s *Storage) FindAllByUserIDAndPeriod(ctx context.Context, userID storage.UserID, from time.Time, to time.Time) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(ctx, selectAllQuery, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(event.ID, event.Title, event.StartAt, event.EndAt, event.Description, event.NotifyAt); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func (s *Storage) HasByUserIDAndPeriod(ctx context.Context, userID storage.UserID, from time.Time, to time.Time) (bool, error) {
	row := s.db.QueryRowContext(ctx, fmt.Sprintf("select exists (%s) as exists", selectAllQuery), userID, from, to)
	if row.Err() != nil {
		return false, row.Err()
	}
	var res struct {
		Exists bool
	}
	if err := row.Scan(&res); err != nil {
		return false, err
	}

	return res.Exists, nil
}

func New(dsn string, maxOpenConns, maxIdleConns int, connMaxLifetime, connMaxIdleTime time.Duration) *Storage {
	return &Storage{
		dsn:             dsn,
		maxOpenConns:    maxOpenConns,
		connMaxLifetime: connMaxLifetime,
		maxIdleConns:    maxIdleConns,
		connMaxIdleTime: connMaxIdleTime,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	if s.db != nil {
		return nil
	}

	db, err := sql.Open("pgx", s.dsn)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(s.maxOpenConns)
	db.SetConnMaxLifetime(s.connMaxLifetime)
	db.SetMaxIdleConns(s.maxIdleConns)
	db.SetConnMaxIdleTime(s.connMaxIdleTime)

	s.db = db

	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) Migrate(ctx context.Context, command string) error {
	if err := s.Connect(ctx); err != nil {
		return fmt.Errorf("db connect: %w", err)
	}
	if err := goose.SetDialect("pgx"); err != nil {
		return fmt.Errorf("set dialect: %w", err)
	}
	goose.SetBaseFS(embedMigrations)

	return goose.Run(command, s.db, "migrations")
}

const saveQuery = `insert into events (id, title, start_at, end_at, description, owner_id, notify_at)
values ($1, $2, $3, $4, $5, $6, $7)
on conflict (id) do update
set title = excluded.title,
	start_at = excluded.start_at,
	end_at = excluded.end_at,
	description = excluded.description,
	notify_at = excluded.notify_at`

const selectQuery = `select id, title, start_at, end_at, description, owner_id, notify_at
from events
where id = $1`

const deleteQuery = `delete from events where id = $1`

const selectAllQuery = `select id, title, start_at, end_at, description, owner_id, notify_at
from events
where owner_id = $1 
  and (start_at between $2 and $3 or end_at between $2 and $3)`
