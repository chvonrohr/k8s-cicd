package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Site is the top most object in the crawling order. Each site has up to n crawls, each of which have n pages.
type Site struct {
	gorm.Model

	Url      string        `json:"url"`
	Interval time.Duration `json:"interval"`
	Crawls   []Crawl       `json:"crawls" gorm:"ForeignKey:SiteID"`
}
