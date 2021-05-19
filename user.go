package delivery_manager

import "time"

type User struct {
	Id           int       `json:"-"                           db:"id"`
	Name         string    `json:"name"     binding:"required" db:"name"`
	Username     string    `json:"username" binding:"required" db:"username"`
	Password     string    `json:"password" binding:"required"`
	PasswordSalt string    `json:"-"                           db:"password_salt"`
	PasswordHash string    `json:"-"                           db:"password_hash"`
	JoinedAt     time.Time `json:"-"                           db:"joined_at"`
}
