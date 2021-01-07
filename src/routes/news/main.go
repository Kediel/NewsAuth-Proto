package newsRoutes

import (
  "github.com/gin-gonic/gin"
  "github.com/z-tech/blue/src/controllers/news"
  "github.com/z-tech/blue/src/middleware"
)

func ApplyToEngine(engine *gin.Engine) {
  engine.POST(
    "/v1/publishArticle",
    middleware.GetConfig,
    newsController.ValidatePublishArticle,
    newsController.PublishArticle,
  )
  engine.POST("/v1/reviseArticle", newsController.ValidateReviseArticle, newsController.ReviseArticle)
  engine.POST("/v1/proveArticle", newsController.ValidateProveArticle, newsController.ProveArticle)
}
