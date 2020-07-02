package model

import "github.com/jinzhu/gorm"

type Crawl struct {
	gorm.Model
	Site   Site   `json:"site"`
	SiteID int    `json:"siteId" gorm:"column:site_id"`
	Pages  []Page `json:"pages" gorm:"ForeignKey:CrawlID"`
}
