package main

import (
  "github.com/gin-gonic/gin"
  "github.com/z-tech/blue/src/routes"
)

func main() {
	engine := gin.Default()
  routes.ApplyAllToEngine(engine)
	engine.Run()
}
