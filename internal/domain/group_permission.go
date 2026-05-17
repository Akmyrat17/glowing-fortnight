package domain

import "time"

type GroupPermission struct {
	ID            int
	Name          string
	PermissionIDs []int
	CreatedAt     time.Time
}
