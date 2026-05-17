package domain

import "time"

type Permission struct {
	ID          int
	Module      string
	Action      string
	Name        string
	Description string
	CreatedAt   time.Time
}
