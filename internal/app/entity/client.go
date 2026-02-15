package entity

import (
	"time"
)

type Client struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CompanyName string    `json:"companyName"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phoneNumber"`
	Address     string    `json:"address"`
	OwnerName   string    `json:"ownerName"`
	IsActive    bool      `json:"isActive"`
	Token       string    `json:"token"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (Client) TableName() string {
	return "client"
}
