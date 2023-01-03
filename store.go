package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func (svc *Service) Create(ctx context.Context) (*Attendee, error) {
	svc.mx.Lock()
	defer svc.mx.Unlock()

	count, err := countAttendees(ctx, svc.db)
	if err != nil {
		return nil, fmt.Errorf("countAttendee(): %w", err)
	}

	var a Attendee
	a.ID = count
	a.CreatedAt = time.Now()
	if err := createAttendee(ctx, svc.db, &a); err != nil {
		return nil, fmt.Errorf("createAttendee(): %w", err)
	}
	return &a, nil
}

type Service struct {
	db *sql.DB
	mx sync.Mutex
}

type Attendee struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func countAttendees(ctx context.Context, db *sql.DB) (int64, error) {
	var n int64
	query := `SELECT COUNT(id) FROM attendees`
	row := db.QueryRowContext(ctx, query)
	err := row.Scan(&n)
	if err != nil {
		return 0, err
	}
	return n + 1, nil
}

func createAttendee(ctx context.Context, db *sql.DB, a *Attendee) error {
	query := `INSERT INTO attendees(id, created_at) VALUES($1, $2)`
	if _, err := db.ExecContext(ctx, query, a.ID, a.CreatedAt); err != nil {
		return err
	}
	return nil
}
