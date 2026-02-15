package entity

import (
	"time"
)

const (
	RoleAdminCode   = 1
	RoleEmployeCode = 2
)

// User represents a user in the system.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	ClientID     uint      `json:"clientId"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	FullName     string    `json:"fullName"`
	Role         uint      `json:"role"`
	Token        string    `json:"token"`
	TokenExpired time.Time `json:"tokenExpired"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
}

func (User) TableName() string {
	return "user"
}
