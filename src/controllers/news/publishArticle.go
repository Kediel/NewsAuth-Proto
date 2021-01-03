package newsController

import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "github.com/go-ozzo/ozzo-validation/v4"

  "github.com/z-tech/blue/src/datalayers/log"
  "github.com/z-tech/blue/src/datalayers/map"
)

type PublishArticleData struct {
  ArticleBody string `json:"ArticleBody"`
  Author string `json:"Author"`
  Dateline string `json:"Dateline"`
}

func (publishArticleData PublishArticleData) Validate() error {
  return validation.ValidateStruct(&publishArticleData,
    validation.Field(&publishArticleData.ArticleBody, validation.Required, validation.Length(1, 20000)),
    validation.Field(&publishArticleData.Author, validation.Required, validation.Length(1, 1000)),
    validation.Field(&publishArticleData.Dateline, validation.Length(0, 1000)),
  )
}

func ValidatePublishArticle(ctx *gin.Context) {
  publishArticleData := PublishArticleData{}
  bindErr := ctx.ShouldBindBodyWith(&publishArticleData, binding.JSON)
  if bindErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
    ctx.Abort()
    return
  }
  validateErr := publishArticleData.Validate()
  if validateErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("one or more properties in request body are not valid: %s", validateErr)})
    ctx.Abort()
    return
  }
  ctx.Set("publishArticleData", publishArticleData)
}

func PublishArticle(ctx *gin.Context) {
  publishArticleData, _ := ctx.Get("publishArticleData")
  leafData, marshalErr := json.Marshal(publishArticleData)
  if marshalErr != nil {
    fmt.Println("error: unable to marshal publishArticleData")
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  leafIndex, treeSize, inclusionProof, rootHash, leafHash, isDup, addLeafErr := logDatalayer.AddLeaf(ctx, leafData)
  if addLeafErr != nil {
    fmt.Println("error: unable to add leaf to log")
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  key := leafHash[:]
  if (isDup == false) {
    addLeafErr := mapDatalayer.AddLeaf(ctx, key, leafData)
    if addLeafErr != nil {
      ctx.JSON(http.StatusInternalServerError, gin.H{})
      ctx.Abort()
      return
    }
  }
  _, mapLeafHash, mapLeafValue, getLeafErr := mapDatalayer.GetLeaf(ctx, key)
  if getLeafErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  ctx.JSON(200, gin.H{
    "ArticleKey": leafHash,
    "ArticleRevisionKey": leafHash,
    "PreviousArticleRevisionKey": nil,
    "InclusionProof": inclusionProof,
    "NonInclusionProof": nil,
    "LogLeafHash": leafHash,
    "LogLeafIndex": leafIndex,
    "LogRootHash": rootHash,
    "LogTreeSize": treeSize,
    "MapLeafHash": mapLeafHash,
    "MapLeafValue": mapLeafValue,
  })
}
