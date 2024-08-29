package entity

import (
	userEntity "github.com/sallescosta/conduit-api/internal/entity/user"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Slug           string           `json:"slug,omitempty"`
	Title          string           `json:"title,omitempty"`
	Description    string           `json:"description,omitempty"`
	Body           string           `json:"body,omitempty"`
	TagList        []string         `json:"tag_list,omitempty"`
	Favorited      bool             `json:"favorited,omitempty"`
	FavoritesCount uint             `json:"favorites_count,omitempty"`
	Author         *userEntity.User `gorm:"embedded;embeddedPrefix:author_"`
}

func NewArticle(
	slug, title, description, body string,
	tagList []string,
	favorited bool,
	author *userEntity.User) (*Article, error) {
	return &Article{
		Slug:           slug,
		Title:          title,
		Description:    description,
		Body:           body,
		TagList:        tagList,
		Favorited:      favorited,
		FavoritesCount: 0,
		Author:         author,
	}, nil
}
