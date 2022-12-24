package app

import "time"

type User struct {
	ID       uint64 `gorm:"primary_key:auto_increment" json:"id"`
	Username string `gorm:"column:username;unique" json:"username"`
	Password string `gorm:"column:password" json:"password"`
}

type Secret struct {
	ID         uint64    `gorm:"primary_key:auto_increment" json:"id"`
	ExternalID string    `gorm:"type:varchar(255)" json:"external_id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	Title      string    `gorm:"column:title;unique" json:"title"`
	Content    string    `gorm:"column:content;unique" json:"content"`
}
