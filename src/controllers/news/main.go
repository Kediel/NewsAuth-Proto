package newsController

import (
  b64 "encoding/base64"
  "crypto/sha256"
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

// TODO(z-tech): what are the fields we actually want?
type PostNewsSchema struct {
  ArticleBody string `json:"ArticleBody"`
  Author string `json:"Author"`
  Dateline string `json:"Dateline"`
}

func (articleIDSchema ArticleIDSchema) Validate() error {
  return validation.ValidateStruct(&articleIDSchema,
    validation.Field(&articleIDSchema.ArticleID, validation.Required, validation.Length(44, 44)),
  )
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
  ctx.Set("postNewsSchema", postNewsSchema)

  validateErr := postNewsSchema.Validate()
  if validateErr != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("one or more properties in request body are not valid: %s", validateErr)})
    ctx.Abort()
    return
  }
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

func PostNews(ctx *gin.Context) {
  postNewsSchema, _ := ctx.Get("postNewsSchema")
  leafData, marshalErr := json.Marshal(postNewsSchema)
  if marshalErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  proof, isDup, getProofErr := logDatalayer.AddLeaf(ctx, leafData)
  if getProofErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  hash := sha256.Sum256(leafData)
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
    ctx.JSON(200, gin.H{"ArticleID": key, "proof": proof, "mapLeaf": mapLeaf, "mapRoot": mapRoot})
    ctx.Abort()
    return
  }

  ctx.JSON(201, gin.H{"ArticleID": key, "proof": proof, "mapLeaf": mapLeaf, "mapRoot": mapRoot})
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

  proof, isDup, getProofErr := logDatalayer.AddLeaf(ctx, leafData)
  if getProofErr != nil {
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
