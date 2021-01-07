package newsRoutes

import (
  "github.com/gin-gonic/gin"
  "github.com/z-tech/blue/src/controllers/news"
)

func ApplyToEngine(engine *gin.Engine) {
  engine.POST(
    "/v1/publishArticle",
    newsController.ValidatePublishArticle,
    newsController.PublishArticle,
  )
  engine.POST(
    "/v1/reviseArticle",
    newsController.ValidateReviseArticle,
    newsController.ReviseArticle,
  )
  engine.POST(
    "/v1/proveArticle",
    newsController.ValidateProveArticle,
    newsController.ProveArticle,
  )
}
