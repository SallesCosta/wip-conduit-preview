package entity

import (
  "github.com/sallescosta/conduit-api/pkg/entity"
)

type Tag struct {
  ID   entity.ID `json:"id"`
  Name string    `json:"name"`
}

func NewTag(tag string) *Tag {
  return &Tag{
    ID:   entity.NewID(),
    Name: tag,
  }
}
