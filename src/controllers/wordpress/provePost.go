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
  _, _, mapAddress, mapID, getConfigErr := envDatalayer.GetConfig()
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
  isExists, _, mapLeafValue, proof, getMapLeafErr := grpcDatalayer.GetMapLeaf(
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

  // 4) potentially return with proof of noninclusion
  if (isExists == false) {
    ctx.JSON(200, gin.H{
      "IsExists": false,
      "IsMostCurrent": false,
      "InclusionProof": nil,
      "NonInclusionProof": proof,
    })
    ctx.Abort()
    return
  }

  // 5) check if most fresh and potentially return
  leafData, marshalErr := json.Marshal(wordpressPost.Data)
  if marshalErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
    ctx.Abort()
    return
  }
  if (reflect.DeepEqual(leafData, mapLeafValue)) {
    ctx.JSON(200, gin.H{
      "IsExists": true,
      "IsMostCurrent": true,
      "InclusionProof": proof,
      "NonInclusionProof": nil,
    })
  }

  // TODO: else fetch log inclusion proof, which in combination with the map
  // noninclusion proof shows the post exists, but is not most current
}
