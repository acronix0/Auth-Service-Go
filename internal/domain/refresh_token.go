package domain

import "time"

type RefreshToken struct {
	ID         int 
	UserID     int
	Token      string 
	DeviceInfo string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}
