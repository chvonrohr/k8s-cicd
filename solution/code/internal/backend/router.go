package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
	"log"
	"strconv"
)

func InitialiseRouter(r *gin.Engine, db *gorm.DB) {

	// create transaction for each request
	r.Use(PersistenceMiddleware(db))

	r.GET("", func(c *gin.Context) {
		c.String(200, "backend works")
	})

	r.POST("/callback/:pageId", func(c *gin.Context) {
		tx := GetTx(c)
		parentId, err := strconv.Atoi(c.Param("pageId"))
		if err != nil {
			_ = c.AbortWithError(500, err)
			return
		}

		var parent model.Page
		tx.First(&parent, parentId)

		var urls []string
		if err := c.BindJSON(&urls); err != nil {
			_ = c.AbortWithError(500, err)
			return
		}

		for _, url := range urls {
			var page model.Page
			tx.Where(model.Page{Site: parent.Site, Url: url}).Attrs(model.Page{State: model.CreatedState}).FirstOrCreate(&page)
			if page.State == model.CreatedState {
				// queue page
				if err := QueuePage(page); err != nil {
					log.Println(err)
					continue
				}
				page.State = model.PendingState
			}
			page.Parents = append(page.Parents, parent)
			tx.Save(&page)
		}
		parent.State = model.CrawledState
		tx.Save(&parent)

		c.Status(200)
	})

	r.GET("/sites", func(c *gin.Context) {
		var sites = make([]model.Site, 1)
		GetTx(c).Find(&sites)
		c.JSON(200, &sites)
		return
	})
	r.POST("/sites", func(c *gin.Context) {
		var site model.Site
		if err := c.BindJSON(&site); err != nil {
			c.AbortWithStatusJSON(500, err)
			return
		}
		GetTx(c).Save(&site)

	})
	r.POST("/sites/:id/crawl", func(c *gin.Context) {
		tx := GetTx(c)
		siteId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			_ = c.AbortWithError(500, err)
			return
		}

		var site model.Site
		tx.First(&site, siteId)

		page := model.Page{
			Site:  site,
			Url:   site.Url,
			State: model.PendingState,
		}
		tx.Create(&page)
		err = QueuePage(page)
		if err != nil {
			FailTx(c)
			c.AbortWithStatusJSON(500, err)
			return
		}
		c.Status(204)
		return

	})

}
