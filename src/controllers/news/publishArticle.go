package newsController

import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "github.com/google/trillian/merkle/rfc6962"

  "github.com/z-tech/blue/src/datalayers"
  "github.com/z-tech/blue/src/types"
)

func ValidatePublishArticle(ctx *gin.Context) {
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

func PublishArticle(ctx *gin.Context) {
  article, _ := ctx.Get("article")
  leafData, _ := json.Marshal(article)
  logAddress, logID, mapAddress, mapID, getConfigErr := datalayers.GetConfig()
  if getConfigErr != nil {
    fmt.Println("error: unable to read config from env %+v\n", getConfigErr)
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  addLogLeafErr := datalayers.AddLogLeaf(ctx, logAddress, logID, leafData)
  if addLogLeafErr != nil {
    fmt.Println("error: unable to add log leaf %+v\n", addLogLeafErr)
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  mapIndex := rfc6962.DefaultHasher.HashLeaf(leafData)
  addMapLeafErr := datalayers.AddMapLeaf(ctx, mapAddress, mapID, mapIndex, leafData)
  if addMapLeafErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  ctx.JSON(200, gin.H{})
}
