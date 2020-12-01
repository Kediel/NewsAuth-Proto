package newsRoutes

import (
  "github.com/gin-gonic/gin"

  "github.com/z-tech/blue/src/controllers/news"
)

func ApplyToEngine(engine *gin.Engine) {
  engine.POST("/v1/news", newsController.PostNews)
}
