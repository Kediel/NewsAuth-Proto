package healthController

import "github.com/gin-gonic/gin"

func Ping(c *gin.Context) {
  // TODO(z-tech): actually check backing services
  c.JSON(200, gin.H{
    "message": "pong",
  })
}
