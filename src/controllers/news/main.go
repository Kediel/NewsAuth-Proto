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

// TODO(z-tech): what are the fields we actually want?
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
  bindErr := ctx.ShouldBindWith(&postNewsSchema, binding.JSON)
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

func PostNews(ctx *gin.Context) {
  postNewsSchema, _ := ctx.Get("postNewsSchema")
  leafData, marshalErr := json.Marshal(postNewsSchema)
  if marshalErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }

  mapDatalayer.AddLeaf(ctx, "HELLOKEY", leafData)

  proof, isDup, getProofErr := logDatalayer.AddLeaf(ctx, leafData)
  if getProofErr != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{})
    ctx.Abort()
    return
  }
  if isDup != false {
    ctx.JSON(200, proof)
    ctx.Abort()
    return
  }

  ctx.JSON(201, proof)
}
