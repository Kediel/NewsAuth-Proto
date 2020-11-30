package healthRoutes

import (
  "github.com/gin-gonic/gin"

  "github.com/z-tech/blue/src/controllers/health"
)

func ApplyToEngine(engine *gin.Engine) {
	engine.GET("/ping", healthController.Ping)
}
