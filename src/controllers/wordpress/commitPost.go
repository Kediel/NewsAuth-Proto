package wordpressController

import (
  "encoding/binary"
  "fmt"
  "net/http"
  "strconv"

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
  logLeafData := []byte(strconv.FormatUint(wordpressPost.ID, 10) + "," + wordpressPost.Data);
  addLogLeafErr := grpcDatalayer.AddLogLeaf(ctx, logAddress, logID, logLeafData)
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
  mapLeafData := []byte(wordpressPost.Data) // level down from logLeafData
  addMapLeafErr := grpcDatalayer.AddMapLeaf(ctx, mapAddress, mapID, mapIndex, mapLeafData)
  if addMapLeafErr != nil {
    fmt.Printf("error: unable to add map leaf %v\n", addMapLeafErr)
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  ctx.JSON(200, gin.H{})
}
