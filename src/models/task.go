package models

import "time"

type Task struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	UserID      uint
	Title       string    `gorm:"type:varchar(255)"`
	Description string    `gorm:"type:text"`
	Status      string    `gorm:"type:varchar(50);default:'pending'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
