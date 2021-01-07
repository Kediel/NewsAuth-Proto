package types

import "github.com/go-ozzo/ozzo-validation/v4"

type Article struct {
  Author string `json:"Author"`
  Body string `json:"Body"`
  Dateline string `json:"Dateline"`
}

func (article Article) Validate() error {
  return validation.ValidateStruct(&article,
    validation.Field(&article.Author, validation.Required, validation.Length(1, 1000)),
    validation.Field(&article.Body, validation.Required, validation.Length(1, 20000)),
    validation.Field(&article.Dateline, validation.Length(0, 1000)),
  )
}
