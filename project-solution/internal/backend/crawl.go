package backend

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/core/internal/model"
)

// crawlSite wraps around crawlSiteWrapped and wraps the passed siteId in a Site struct
func crawlSite(tx *gorm.DB, siteId uint) error {
	var c model.Crawl
	c.SiteID = siteId
	return crawlSiteWrapped(tx, c)
}

// crawlSiteWrapped creates a site on the current database
// transaction and then queues it on the configured rabbitmq queue.
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
