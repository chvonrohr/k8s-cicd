package model

import "github.com/jinzhu/gorm"

type Crawl struct {
	gorm.Model
	Site   Site   `json:"site"`
	SiteID uint   `json:"siteId" gorm:"column:site_id"`
	Pages  []Page `json:"pages" gorm:"ForeignKey:CrawlID"`
}
