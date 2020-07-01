package model

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	Site       Site   `json:"site" gorm:"ForeignKey:id;References:id"`
	ParentID   int    `json:"parentId"`
	ParentType string `json:"parentType"`
	Url        string `json:"url"`
	Children   []Page `json:"-" gorm:"polymorphic:Parent;"`
}
