package types

import "github.com/go-ozzo/ozzo-validation/v4"

type WordpressPost struct {
  ID uint64 `json:"ID"`
  Data string `json:"Data"`
}

func (wordpressPost WordpressPost) Validate() error {
  return validation.ValidateStruct(&wordpressPost,
    validation.Field(&wordpressPost.ID, validation.Required),
    validation.Field(&wordpressPost.Data, validation.Required),
  )
}
