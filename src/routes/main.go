package routes

import (
  "github.com/gin-gonic/gin"

  "github.com/z-tech/blue/src/routes/health"
  "github.com/z-tech/blue/src/routes/news"
)

func ApplyAllToEngine(engine *gin.Engine) {
  healthRoutes.ApplyToEngine(engine)
  newsRoutes.ApplyToEngine(engine)
}
