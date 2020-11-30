package newsController

import (
  "log"

  "github.com/gin-gonic/gin"
  "github.com/go-ozzo/ozzo-validation/v4"
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
  bindErr := c.BindJSON(&postNewsSchema)
  if bindErr != nil {
    log.Printf("warn: unable to parse request body %+v", bindErr)
    c.AbortWithStatusJSON(400, gin.H{"error": "unable to parse request body"})
    return
  }

  validateErr := postNewsSchema.Validate()
  if validateErr != nil {
    log.Printf("warn: unable to validate request body %+v", validateErr)
    c.AbortWithStatusJSON(400, gin.H{"error": "unable to validate request body properties"})
    return
  }

  c.JSON(201, gin.H{
    "message": "pong",
  })
}
