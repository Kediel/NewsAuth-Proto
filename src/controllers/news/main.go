package newsController

import (
  "encoding/json"
  "fmt"
  "log"

  "github.com/gin-gonic/gin"
  "github.com/go-ozzo/ozzo-validation/v4"

  "github.com/z-tech/blue/src/datalayers/log"
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

func PostNews(c *gin.Context) {
  postNewsSchema := PostNewsSchema{}
  err := c.BindJSON(&postNewsSchema)
  if err != nil {
    log.Printf("warn: unable to parse request body %+v", err)
    c.AbortWithStatusJSON(400, gin.H{"error": "unable to parse request body"})
    return
  }

  err1 := postNewsSchema.Validate()
  if err1 != nil {
    c.AbortWithStatusJSON(400, gin.H{"error": fmt.Sprintf("one or more properties in request body are not valid: %s", err1)})
    return
  }

  leafData, err2 := json.Marshal(postNewsSchema)
  proof, err3 := logDatalayer.AddLeaf(leafData)
  if err2 != nil || err3 != nil {
    c.AbortWithStatusJSON(500, gin.H{"error": "unexpected error"})
    return
  }

  c.JSON(201, gin.H{
    "proof": proof,
  })
}
