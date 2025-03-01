package models

type Message struct {
	ID uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id"`
	Content string `json:"content"`
	Reply string `json:"reply"`
}