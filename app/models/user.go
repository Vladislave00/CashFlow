package models

type User struct {
	ID       int    `json:"user_id" db:"user_id"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}
