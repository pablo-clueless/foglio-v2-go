package middlewares

import (
	"foglio/v2/src/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FileMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		form, err := ctx.MultipartForm()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Unable to parse multipart form",
			})
			return
		}

		files := form.File["images"]
		if len(files) == 0 {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No files found",
			})
			return
		}

		for _, header := range files {
			if header.Size > int64(config.AppConfig.MaxFileSize) {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"message": "File size exceeds 5MB limit",
				})
				return
			}
		}

		ctx.Next()
	}
}
