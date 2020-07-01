package model

import "gorm.io/gorm"

type PageState string

var (
	PendingState PageState = "pending"
	CrawledState PageState = "crawled"
	CreatedState PageState = "created"
)

type Page struct {
	gorm.Model
	Site    Site      `json:"site" gorm:"index;ForeignKey:id;References:id"`
	Url     string    `json:"url" gorm:"index"`
	State   PageState `json:"state"`
	Parents []Page    `json:"-" gorm:"many2many:page_ref_table;association_jointable_foreignkey:parent_id"`
}
