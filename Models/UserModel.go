package Models

import (
  "time"
)

type User struct {
	ID uint							`gorm:"primaryKey"`
	Email string
	FirstName string
	LastName string
	SlackID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *User) TableName() string {
  return "users"
}
