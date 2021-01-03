package newsController

import (
  b64 "encoding/base64"
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/gin-gonic/gin"
  "github.com/gin-gonic/gin/binding"
  "github.com/go-ozzo/ozzo-validation/v4"

  "github.com/z-tech/blue/src/datalayers/log"
  "github.com/z-tech/blue/src/datalayers/map"
)

type ArticleIDSchema struct {
  ArticleID string `json:"ArticleID"`
}

func (articleIDSchema ArticleIDSchema) Validate() error {
  return validation.ValidateStruct(&articleIDSchema,
    validation.Field(&articleIDSchema.ArticleID, validation.Required, validation.Length(44, 44)),
  )
}

func ValidatePostNewsRevision(ctx *gin.Context) {
  ValidatePostNews(ctx)

  articleIDSchema := ArticleIDSchema{}
  bindErr := ctx.ShouldBindBodyWith(&articleIDSchema, binding.JSON)
  if bindErr != nil {
    fmt.Printf("HELLLO 4 %+v\n", bindErr)
    ctx.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
    ctx.Abort()
    return
  }

  validateErr := articleIDSchema.Validate()
  if validateErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("one or more properties in request body are not valid: %s", validateErr)})
    ctx.Abort()
    return
  }
  ctx.Set("articleID", articleIDSchema.ArticleID)
}

func PostNewsRevision(ctx *gin.Context) {
  articleIDRaw, _ := ctx.Get("articleID")
  articleID, isString := articleIDRaw.(string)
  if isString == false {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }
  fmt.Printf("HELLO 1 %+v\n", articleID)

  postNewsSchema, _ := ctx.Get("postNewsSchema")
  leafData, marshalErr := json.Marshal(postNewsSchema)
  if marshalErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  _, _, proof, _, _, isDup, addLeafErr := logDatalayer.AddLeaf(ctx, leafData)
  if addLeafErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  fmt.Printf("HLLOO %+v %v %v\n", proof, isDup, articleID)

  uDec, _ := b64.URLEncoding.DecodeString(articleID)
  hash := []byte(uDec)
  key := hash[:]
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

  if isDup == true {
    ctx.JSON(200, gin.H{"proof": proof, "mapLeaf": mapLeaf, "mapRoot": mapRoot})
    ctx.Abort()
    return
  }

  ctx.JSON(201, gin.H{"proof": proof, "mapLeaf": mapLeaf, "mapRoot": mapRoot})
}
