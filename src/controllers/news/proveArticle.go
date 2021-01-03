package newsController

import (
  b64 "encoding/base64"
  // "encoding/json"
  // "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"

  // "github.com/z-tech/blue/src/datalayers/log"
  "github.com/z-tech/blue/src/datalayers/map"
)

type ProveArticleData struct {
  ArticleKey string `json:"ArticleKey"`
}

func ValidateProveArticle(ctx *gin.Context) {
  // ValidatePublishArticle(ctx)
}

func ProveArticle(ctx *gin.Context) {
  // 1) validate the articleKey
  proveArticleData := ProveArticleData{}
  bindErr := ctx.ShouldBindBodyWith(&proveArticleData, binding.JSON)
  if bindErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
    ctx.Abort()
    return
  }
  articleKey, decodeArticleKeyErr := b64.URLEncoding.DecodeString(proveArticleData.ArticleKey)
  if decodeArticleKeyErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": decodeArticleKeyErr.Error()})
    ctx.Abort()
    return
  }
  if len(articleKey) != 32 {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "width of ArticleKey must be 32"})
    ctx.Abort()
    return
  }

  // 2) get the latest revision of this article from the map
  isExists, mapLeafHash, mapLeafValue, getLeafErr := mapDatalayer.GetLeaf(ctx, articleKey)
  if getLeafErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }
  // TODO: prove non-inclusion
  if isExists != true {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": "ArticleID does not exist"})
    ctx.Abort()
    return
  }

  




  ctx.JSON(200, gin.H{
    // "ArticleKey": leafHash,
    // "ArticleRevisionKey": leafHash,
    // "PreviousArticleRevisionKey": nil,
    // "InclusionProof": inclusionProof,
    // "NonInclusionProof": nil,
    // "LogLeafHash": leafHash,
    // "LogLeafIndex": leafIndex,
    // "LogRootHash": rootHash,
    // "LogTreeSize": treeSize,
    "MapLeafHash": mapLeafHash,
    "MapLeafValue": mapLeafValue,
  })
}
