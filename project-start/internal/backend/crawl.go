package backend

import (
	"gitlab.com/letsboot/core/kubernetes-course/project-vision/internal/model"
	"gorm.io/gorm"
)

// crawlSite wraps around crawlSiteWrapped and wraps the passed siteId in a Site struct
func crawlSite(tx *gorm.DB, siteId uint) (model.Crawl, error) {
	var c model.Crawl
	c.SiteID = siteId
	return crawlSiteWrapped(tx, c)
}

// crawlSiteWrapped creates a site on the current database
// transaction and then queues it on the configured rabbitmq queue.
func crawlSiteWrapped(tx *gorm.DB, c model.Crawl) (model.Crawl, error) {
	tx.Save(&c)
	var site model.Site
	tx.First(&site, c.SiteID)
	page := model.Page{
		Crawl: c,
		Url:   site.Url,
		State: model.PendingState,
	}
	tx.Create(&page)
	return c, QueuePage(page)
}
