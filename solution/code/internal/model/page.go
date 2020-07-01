package model

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	SiteId int    `json:"siteId"`
	Url    string `json:"url"`
}
