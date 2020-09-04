package model

import (
	"gorm.io/gorm"
)

// Crawl represents a single crawl. A crawl contains n Pages and a parent Site.
// It is an abstract object to help us manage multiple "Crawls" for a page while
// maintaining the RESTful principles.
type Crawl struct {
	gorm.Model
	Site   Site   `json:"site"`
	SiteID uint   `json:"siteId" gorm:"column:site_id"`
	Pages  []Page `json:"pages" gorm:"ForeignKey:CrawlID"`
}
