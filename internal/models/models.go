package models

import (
	"time"
)

//Model : a copy of gorm.Model with json annotations
type Model struct {
	ID        uint       `gorm:"primary_key;column:id" json:"id"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type Cluster struct {
	Model
	Name     string `gorm:"not null" json:"name"`
	Provider string `json:"provider"`
	Token    string `json:"token"`
}
