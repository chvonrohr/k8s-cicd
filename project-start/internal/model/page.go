package model

import "gorm.io/gorm"

// PageState is the state of a page, as determined by its current crawl status.
type PageState string

var (
	PendingState PageState = "pending"
	CrawledState PageState = "crawled"
	CreatedState PageState = "created"
	ErroredState PageState = "errored"
)

// Page represents a single url. It contains a status code, a content type and a state.
// It can also have n parent Pages, forming a tree.
type Page struct {
	gorm.Model
	StatusCode    int       `json:"statusCode"`
	ContentType   string    `json:"contentType"`
	Crawl         Crawl     `json:"crawl"`
	CrawlID       int       `json:"crawlId" gorm:"column:crawl_id"`
	Url           string    `json:"url" gorm:"index:url"`
	State         PageState `json:"state"`
	Parents       []Page    `json:"-" gorm:"many2many:page_ref_table;association_jointable_foreignkey:parent_id"`
	FileAvailable bool      `json:"fileAvailable"`
}
