package main

import "github.com/gin-gonic/gin"
import "github.com/z-tech/blue/src/routes"

func main() {
  engine := gin.Default()
  routes.ApplyAllToEngine(engine)
  engine.Run()
}
