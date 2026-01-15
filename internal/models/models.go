package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	Role      string    `json:"role" binding:"required,oneof=elderly family"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CheckIn represents a daily wellness check-in
type CheckIn struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id" binding:"required"`
	CheckInType  string    `json:"check_in_type" binding:"required,oneof=manual passive"`
	StepCount    *int      `json:"step_count,omitempty"`
	BatteryLevel *int      `json:"battery_level,omitempty"`
	CheckedAt    time.Time `json:"checked_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// FamilyConnection represents the relationship between elderly and family users
type FamilyConnection struct {
	ID            int       `json:"id"`
	ElderlyUserID int       `json:"elderly_user_id" binding:"required"`
	FamilyUserID  int       `json:"family_user_id" binding:"required"`
	Relationship  string    `json:"relationship"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateUserRequest is the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Role  string `json:"role" binding:"required,oneof=elderly family"`
}

// CreateCheckInRequest is the request body for creating a check-in
type CreateCheckInRequest struct {
	CheckInType  string `json:"check_in_type" binding:"required,oneof=manual passive"`
	StepCount    *int   `json:"step_count,omitempty"`
	BatteryLevel *int   `json:"battery_level,omitempty"`
}

// CheckInResponse includes user information with the check-in
type CheckInResponse struct {
	CheckIn
	User User `json:"user"`
}
