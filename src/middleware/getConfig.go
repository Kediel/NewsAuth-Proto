package middleware

import (
  "os"
  "fmt"
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
)

func GetConfig(ctx *gin.Context) {
  LOG_ADDRESS := os.Getenv("LOG_ADDRESS")
  LOG_ID, logIDConvErr := strconv.ParseInt(os.Getenv("LOG_ID"), 10, 64)
  if logIDConvErr != nil {
    fmt.Printf("error: unable read log id from environment %+v\n", logIDConvErr)
    ctx.JSON(http.StatusBadRequest, gin.H{"error": logIDConvErr.Error()})
    ctx.Abort()
    return
  }
  MAP_ADDRESS := os.Getenv("MAP_ADDRESS")
  MAP_ID, mapIDConvErr := strconv.ParseInt(os.Getenv("MAP_ID"), 10, 64)
  if mapIDConvErr != nil {
    fmt.Printf("error: unable read map id from environment %+v\n", mapIDConvErr)
    ctx.JSON(http.StatusBadRequest, gin.H{"error": mapIDConvErr.Error()})
    ctx.Abort()
    return
  }
  ctx.Set("LOG_ADDRESS", LOG_ADDRESS)
  ctx.Set("LOG_ID", LOG_ID)
  ctx.Set("MAP_ADDRESS", MAP_ADDRESS)
  ctx.Set("MAP_ID", MAP_ID)
}
