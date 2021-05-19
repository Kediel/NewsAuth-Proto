package wordpressController

import (
  "encoding/binary"
  "encoding/json"
  "fmt"
  "net/http"
  "reflect"

  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "github.com/google/trillian/merkle/rfc6962"

  "github.com/z-tech/blue/src/datalayers/env"
  "github.com/z-tech/blue/src/datalayers/grpc"
  "github.com/z-tech/blue/src/types"
)

func ProvePost(ctx *gin.Context) {
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

  // 3) get proof from map
  buf := make([]byte, 8)
  binary.LittleEndian.PutUint64(buf, wordpressPost.ID)
  mapIndex := rfc6962.DefaultHasher.HashLeaf(buf)
  isIndexSetInMap, _, mapLeafValue, mapProof, getMapLeafErr := grpcDatalayer.GetMapLeaf(
    ctx,
    mapAddress,
    mapID,
    mapIndex,
  )
  if getMapLeafErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  // 4) get proof from log
  leafData, marshalErr := json.Marshal(wordpressPost)
  if marshalErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
    ctx.Abort()
    return
  }
  leafIndex, _, logProof, _, _, getLogLeafErr := grpcDatalayer.GetLogLeaf(
    ctx,
    logAddress,
    logID,
    leafData,
  )
  if getLogLeafErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  // 5) returns

  // 5a) this index is not set in the map, so, not included
  if (isIndexSetInMap == false) {
    ctx.JSON(200, gin.H{
      "IsIncluded": false,
      "IsMostCurrent": false,
      "LogInclusionProof": logProof,
      "MapInclusionProof": nil,
      "MapNonInclusionProof": mapProof,
    })
    ctx.Abort()
    return
  }

  // 5b) the value at this index of the map is this value, so, fresh and included
  if (reflect.DeepEqual(wordpressPost.Data, string(mapLeafValue))) {
    ctx.JSON(200, gin.H{
      "IsIncluded": true,
      "IsMostCurrent": true,
      "LogInclusionProof": logProof,
      "MapInclusionProof": mapProof,
      "MapNonInclusionProof": nil,
    })
    return
  }

  // 5c) this index is set, but not with this value, and never with this value, perhaps a spoof
  if (leafIndex == -1) {
    ctx.JSON(200, gin.H{
      "IsIncluded": false,
      "IsMostCurrent": false,
      "LogInclusionProof": logProof,
      "MapInclusionProof": nil,
      "MapNonInclusionProof": mapProof,
    })
    return
  }

  // 5d) this is a genuine version, but is outdated
  ctx.JSON(200, gin.H{
    "IsIncluded": true,
    "IsMostCurrent": false,
    "LogInclusionProof": logProof,
    "MapInclusionProof": nil,
    "MapNonInclusionProof": mapProof,
  })
}
