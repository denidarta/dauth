package core

import "time"

// Need to add : device fingerprint

type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Product struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type UserProduct struct {
	UserID    string
	ProductID string
	Role      string
	CreatedAt time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	ProductID string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Invitation struct {
	ID         string
	Email      string
	ProductID  string
	Role       string
	TokenHash  string
	ExpiresAt  time.Time
	AcceptedAt *time.Time
	CreatedAt  time.Time
}
