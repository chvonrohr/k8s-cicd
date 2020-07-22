package backend

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
)

func crawlSite(tx *gorm.DB, siteId uint) error {
	var c model.Crawl
	c.SiteID = siteId
	return crawlSiteWrapped(tx, c)
}

func crawlSiteWrapped(tx *gorm.DB, c model.Crawl) error {
	tx.Save(&c)
	var site model.Site
	tx.First(&site, c.SiteID)
	page := model.Page{
		Crawl: c,
		Url:   site.Url,
		State: model.PendingState,
	}
	tx.Create(&page)
	return QueuePage(page)
}
