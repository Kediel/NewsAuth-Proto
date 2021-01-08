package wordpressRoutes

import (
  "github.com/gin-gonic/gin"

  "github.com/z-tech/blue/src/controllers/wordpress"
)

func ApplyToEngine(engine *gin.Engine) {
  engine.POST("/v1/commitWordpressPost", wordpressController.CommitPost)
  engine.POST("/v1/proveWordpressPost", wordpressController.ProvePost)
}
