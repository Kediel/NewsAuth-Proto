package wordpressController

import (
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"

  "github.com/z-tech/blue/src/datalayers/env"
  "github.com/z-tech/blue/src/datalayers/grpc"
)

func TreeRoots(ctx *gin.Context) {
  // 1) get some config
  logAddress, logID, mapAddress, mapID, getConfigErr := envDatalayer.GetConfig()
  if getConfigErr != nil {
    fmt.Printf("error: unable to read config from env %v\n", getConfigErr)
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  // 2) get map root
  mapRoot, _ := grpcDatalayer.GetMapRoot(ctx, mapAddress, mapID)

  // 3) get log root
  logRoot, _ := grpcDatalayer.GetLogRoot(ctx, logAddress, logID)

  // 4) return
  ctx.JSON(200, gin.H{
    "MapRoot": mapRoot,
    "LogRoot": logRoot,
  })
}
