package model

import "github.com/jinzhu/gorm"

type PageState string

var (
	PendingState PageState = "pending"
	CrawledState PageState = "crawled"
	CreatedState PageState = "created"
)

type Page struct {
	gorm.Model
	Site    Site      `json:"site"`
	SiteID  int       `json:"siteId" gorm:"column:site_id"`
	Url     string    `json:"url" gorm:"index:url"`
	State   PageState `json:"state"`
	Parents []Page    `json:"-" gorm:"many2many:page_ref_table;association_jointable_foreignkey:parent_id"`
}
