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

type PostNewsSchema struct {
  ArticleBody string `json:"ArticleBody"`
  Author string `json:"Author"`
  Dateline string `json:"Dateline"`
}

func (postNewsSchema PostNewsSchema) Validate() error {
  return validation.ValidateStruct(&postNewsSchema,
    validation.Field(&postNewsSchema.ArticleBody, validation.Required, validation.Length(1, 20000)),
    validation.Field(&postNewsSchema.Author, validation.Required, validation.Length(1, 1000)),
    validation.Field(&postNewsSchema.Dateline, validation.Length(0, 1000)),
  )
}

func ValidatePostNews(ctx *gin.Context) {
  postNewsSchema := PostNewsSchema{}
  bindErr := ctx.ShouldBindBodyWith(&postNewsSchema, binding.JSON)
  if bindErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
    ctx.Abort()
    return
  }
  validateErr := postNewsSchema.Validate()
  if validateErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("one or more properties in request body are not valid: %s", validateErr)})
    ctx.Abort()
    return
  }
  ctx.Set("postNewsSchema", postNewsSchema)
}

func PostNews(ctx *gin.Context) {
  postNewsSchema, _ := ctx.Get("postNewsSchema")
  leafData, marshalErr := json.Marshal(postNewsSchema)
  if marshalErr != nil {
    fmt.Println("error: unable to marshal postNewsSchema")
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  leafIndex, treeSize, proof, rootHash, leafHash, isDup, addLeafErr := logDatalayer.AddLeaf(ctx, leafData)
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
  mapLeaf, mapRoot, getLeafErr := mapDatalayer.GetLeaf(ctx, key)
  if getLeafErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  ctx.JSON(200, gin.H{"LogProof": proof, "LogLeafIndex": leafIndex, "LogTreeSize": treeSize, "LogRootHash": rootHash, "LogLeafHash": leafHash, "mapLeaf": mapLeaf, "mapRoot": mapRoot})
}
