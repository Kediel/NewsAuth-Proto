package routes

import (
  "github.com/gin-gonic/gin"

  "github.com/z-tech/blue/src/routes/health"
  "github.com/z-tech/blue/src/routes/wordpress"
)

func ApplyAllToEngine(engine *gin.Engine) {
  healthRoutes.ApplyToEngine(engine)
  wordpressRoutes.ApplyToEngine(engine)
}
