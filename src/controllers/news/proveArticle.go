package newsController

import (
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

func ValidateProveArticle(ctx *gin.Context) {
  article := types.Article{}
  bindErr := ctx.ShouldBindBodyWith(&article, binding.JSON)
  if bindErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
    ctx.Abort()
    return
  }
  validateErr := article.Validate()
  if validateErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("one or more properties in request body are not valid: %s", validateErr)})
    ctx.Abort()
    return
  }
  ctx.Set("article", article)
}

func ProveArticle(ctx *gin.Context) {
  article, _ := ctx.Get("article")
  leafData, _ := json.Marshal(article)
  _, _, mapAddress, mapID, getConfigErr := envDatalayer.GetConfig()
  if getConfigErr != nil {
    fmt.Println("error: unable to read config from env %+v\n", getConfigErr)
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  mapIndex := rfc6962.DefaultHasher.HashLeaf(leafData)
  isExists, mapLeafHash, mapLeafValue, proof, getMapLeafErr := grpcDatalayer.GetMapLeaf(
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

  var inclusionProof [][]byte
  var nonInclusionProof [][]byte
  if (isExists == true) {
    inclusionProof = proof
  } else {
    nonInclusionProof = proof
  }

  ctx.JSON(200, gin.H{
    // "ArticleKey": leafHash,
    // "ArticleRevisionKey": leafHash,
    // "PreviousArticleRevisionKey": nil,
    "InclusionProof": inclusionProof,
    "NonInclusionProof": nonInclusionProof,
    // "LogLeafHash": leafHash,
    // "LogLeafIndex": leafIndex,
    // "LogRootHash": rootHash,
    // "LogTreeSize": treeSize,
    "MapLeafHash": mapLeafHash,
    "MapLeafValue": mapLeafValue,
  })
}
