package models

import (
	"time"
)

type AccountType string
type AccountPermission string
type AccountStatus string

const (
	UserRole  AccountPermission = "user"
	AdminRole AccountPermission = "admin"
)

const (
	ActiveAccount    AccountStatus = "Active"
	SuspendedAccount AccountStatus = "Suspended"
	InactiveAccount  AccountStatus = "Inactive"
)

type User struct {
	ID            uint              `gorm:"primarykey"`
	Email         string            `json:"email"`
	Password      string            `json:"password"`
	LastLogin     string            `json:"last_login"`
	IP            string            `json:"ip"`
	UserId        string            `json:"user_id" validate:"required"`
	Role          AccountPermission `json:"role" validate:"required"`
	EmailVerified bool              `json:"email_verified" validate:"required"`
	Country       string            `json:"country"`
	PhoneNumber   string            `json:"phone_number"`
	Status        AccountStatus     `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	FirstName     string            `json:"first_name"`
	LastName      string            `json:"last_name"`
}
