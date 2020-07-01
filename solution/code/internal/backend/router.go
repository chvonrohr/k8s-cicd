package backend

import "github.com/gin-gonic/gin"

func InitialiseRouter(r *gin.Engine) {

	r.GET("", func(context *gin.Context) {
		context.String(200, "backend works")
	})

}
