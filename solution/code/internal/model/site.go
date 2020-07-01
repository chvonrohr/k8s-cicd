package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Site struct {
	gorm.Model
	Url      string        `json:"url"`
	Interval time.Duration `json:"interval"`
}
