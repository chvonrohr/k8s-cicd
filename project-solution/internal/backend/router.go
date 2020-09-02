package backend

import (
	"log"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/model"
	"gitlab.com/letsboot/core/kubernetes-course/project-solution/internal/sdk"
	"gorm.io/gorm"
)

// InitialiseRouter sets up a new routing instance and configures paths as well as requests.
func InitialiseRouter(r *gin.Engine, db *gorm.DB) {

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowHeaders:     []string{"content-type", "Origin"},
		MaxAge:           12 * time.Hour,
	}))

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

		var response sdk.PageResponse
		if err := c.BindJSON(&response); err != nil {
			_ = c.AbortWithError(500, err)
			return
		}

		if response.Ok {
			parent.State = model.CrawledState
		} else {
			parent.State = model.ErroredState
		}
		for _, url := range response.Urls {
			var page model.Page
			tx.Where(model.Page{CrawlID: parent.CrawlID, Url: url}).Attrs(model.Page{State: model.CreatedState, CrawlID: parent.CrawlID}).FirstOrCreate(&page)
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
		parent.StatusCode = response.StatusCode
		parent.ContentType = response.ContentType
		tx.Save(&parent)

		c.Status(200)
	})

	r.GET("/sites", func(c *gin.Context) {
		var sites []model.Site
		GetTx(c).Find(&sites)
		c.JSON(200, &sites)
		return
	})
	r.GET("/crawls", func(c *gin.Context) {
		var crawls []model.Crawl
		tx := GetTx(c)
		if siteQuery := c.Query("site"); siteQuery != "" {
			tx = tx.Where("site_id = ?", siteQuery)
		}
		tx.Find(&crawls)
		c.JSON(200, &crawls)
		return
	})
	r.GET("/pages", func(c *gin.Context) {
		var pages []model.Page
		tx := GetTx(c)
		if crawlQuery := c.Query("crawl"); crawlQuery != "" {
			tx = tx.Where("crawl_id = ?", crawlQuery)
		}
		tx.Find(&pages)
		c.JSON(200, &pages)
		return
	})

	r.POST("/sites", func(c *gin.Context) {
		var site model.Site
		if err := c.BindJSON(&site); err != nil {
			c.AbortWithStatusJSON(500, err)
			return
		}
		GetTx(c).Save(&site)
		c.JSON(200, site)
		return
	})
	r.POST("/crawls", func(c *gin.Context) {
		tx := GetTx(c)
		var crawl model.Crawl
		if err := c.BindJSON(&crawl); err != nil {
			c.AbortWithStatusJSON(500, err)
			return
		}
		crawl, err := crawlSiteWrapped(tx, crawl)
		if err != nil {
			FailTx(c)
			c.AbortWithStatusJSON(500, err)
			return
		}
		c.Status(204)
		return

	})

	r.POST("/schedule", func(c *gin.Context) {
		tx := GetTx(c)
		var sites []model.Site
		tx.Preload("Crawls", func(db *gorm.DB) *gorm.DB {
			return db.Order("crawls.created_at DESC")
		}).Find(&sites)
		response := sdk.SchedulerResponse{}
		for _, site := range sites {
			if len(site.Crawls) > 0 && time.Since(site.Crawls[0].CreatedAt) < site.Interval {
				continue
			}
			// create crawl for site
			crawl, err := crawlSite(tx, site.ID)
			if err != nil {
				FailTx(c)
				c.AbortWithStatusJSON(500, err)
				return
			}
			response.ScheduledCount++
			response.ScheduledCrawls = append(response.ScheduledCrawls, crawl)
		}
		c.JSON(200, response)
		return

	})
}
