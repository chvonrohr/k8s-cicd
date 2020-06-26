package model

type Page struct {
	Id     int    `json:"id"`
	SiteId int    `json:"siteId"`
	Url    string `json:"url"`
}
