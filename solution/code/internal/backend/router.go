package backend

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/letsboot/core/kubernetes-course/solution/code/core/internal/model"
	"gorm.io/gorm"
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
		pageId, err := strconv.Atoi(c.Param("pageId"))
		if err != nil {
			_ = c.AbortWithError(500, err)
			return
		}

		var page model.Page
		tx.First(&page, pageId)

	})

}
