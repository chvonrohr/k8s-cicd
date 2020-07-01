package backend

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
	"gorm.io/gorm"
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

}
