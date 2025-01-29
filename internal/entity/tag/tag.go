package entity

import (
	"github.com/sallescosta/conduit-api/pkg/entity"
)

type Tag struct {
	ID   entity.ID `json:"id"`
	Name string    `json:"tag"`
}

func NewTag(tag string) *Tag {
	return &Tag{
		ID:   entity.NewID(),
		Name: tag,
	}
}
