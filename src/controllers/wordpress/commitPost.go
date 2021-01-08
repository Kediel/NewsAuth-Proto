package wordpressController

import (
  "encoding/binary"
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "github.com/google/trillian/merkle/rfc6962"

  "github.com/z-tech/blue/src/datalayers/env"
  "github.com/z-tech/blue/src/datalayers/grpc"
  "github.com/z-tech/blue/src/types"
)

func CommitPost(ctx *gin.Context) {
  // 1) validate some stuff
  wordpressPost := types.WordpressPost{}
  bindErr := ctx.ShouldBindBodyWith(&wordpressPost, binding.JSON)
  if bindErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
    ctx.Abort()
    return
  }
  validateErr := wordpressPost.Validate()
  if validateErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("one or more properties in request body are not valid: %s", validateErr)})
    ctx.Abort()
    return
  }

  // 2) get some config
  logAddress, logID, mapAddress, mapID, getConfigErr := envDatalayer.GetConfig()
  if getConfigErr != nil {
    fmt.Printf("error: unable to read config from env %v\n", getConfigErr)
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  // 3) add leaf to log
  leafData, marshalErr := json.Marshal(wordpressPost.Data)
  if marshalErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
    ctx.Abort()
    return
  }
  addLogLeafErr := grpcDatalayer.AddLogLeaf(ctx, logAddress, logID, leafData)
  if addLogLeafErr != nil {
    fmt.Printf("error: unable to add log leaf %v\n", addLogLeafErr)
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  // 4) add leaf to map
  buf := make([]byte, 8)
  binary.LittleEndian.PutUint64(buf, wordpressPost.ID)
  mapIndex := rfc6962.DefaultHasher.HashLeaf(buf)
  addMapLeafErr := grpcDatalayer.AddMapLeaf(ctx, mapAddress, mapID, mapIndex, leafData)
  if addMapLeafErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  ctx.JSON(200, gin.H{})
}
