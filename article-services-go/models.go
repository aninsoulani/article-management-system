package main

import (
	"time"
)

type Post struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    Title       string    `gorm:"size:200" json:"title" validate:"required,min=20"`
    Content     string    `gorm:"type:text" json:"content" validate:"required,min=200"`
    Category    string    `gorm:"size:100" json:"category" validate:"required,min=3"`
    Status      string    `gorm:"size:100" json:"status" validate:"required,oneof=publish draft trash"`
    CreatedDate time.Time `gorm:"autoCreateTime" json:"created_date"`
    UpdatedDate time.Time `gorm:"autoUpdateTime" json:"updated_date"`
}