package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type PermissionSeed struct {
	Module      string `json:"module"`
	Action      string `json:"action"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GroupPermissionSeed struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type AdminUserSeed struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	GroupPermission string `json:"group_permission"`
}

type SeedData struct {
	Permissions      []PermissionSeed      `json:"permissions"`
	GroupPermissions []GroupPermissionSeed `json:"group_permissions"`
	AdminUser        AdminUserSeed         `json:"admin_user"`
}

func LoadSeedData(path string) (*SeedData, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read seed data: %w", err)
	}

	var data SeedData
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("failed to parse seed data: %w", err)
	}

	return &data, nil
}
