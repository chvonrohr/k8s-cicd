package model

import (
	"gorm.io/gorm"
	"time"
)

type Site struct {
	gorm.Model
	Url      string        `json:"url"`
	Interval time.Duration `json:"interval"`
	Children []Page        `json:"-" gorm:"polymorphic:Parent;"`
}
