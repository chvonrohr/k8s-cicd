package model

import "github.com/jinzhu/gorm"

type PageState string

var (
	PendingState PageState = "pending"
	CrawledState PageState = "crawled"
	CreatedState PageState = "created"
	ErroredState PageState = "errored"
)

type Page struct {
	gorm.Model
	StatusCode  int       `json:"statusCode"`
	ContentType string    `json:"contentType"`
	Crawl       Crawl     `json:"crawl"`
	CrawlID     int       `json:"crawlId" gorm:"column:crawl_id"`
	Url         string    `json:"url" gorm:"index:url"`
	State       PageState `json:"state"`
	Parents     []Page    `json:"-" gorm:"many2many:page_ref_table;association_jointable_foreignkey:parent_id"`
}
