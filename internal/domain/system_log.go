package domain

import "time"

type SystemLog struct {
	ID        string
	Level     string
	Type      string
	Message   string
	Context   map[string]interface{}
	CreatedAt time.Time
}
