package entity

import (
	"github.com/gosimple/slug"
	"github.com/sallescosta/conduit-api/pkg/entity"
	"time"
)

type Article struct {
	ID             entity.ID `json:"id"`
	Slug           string    `json:"slug,omitempty"`
	Title          string    `json:"title,omitempty"`
	Description    string    `json:"description,omitempty"`
	Body           string    `json:"body,omitempty"`
	TagList        []string  `json:"tag_list,omitempty"`
	Favorited      bool      `json:"favorited,omitempty"`
	FavoritesCount uint      `json:"favorites_count,omitempty"`
	CreatedAt      string    `json:"created_at,omitempty"`
	UpdatedAt      string    `json:"updated_at,omitempty"`
	AuthorID       string    `gorm:"embedded;embeddedPrefix:author_"`
}

type AllArticlesOutput struct {
	Articles      []Article `json:"articles"`
	ArticlesCount int       `json:"articlesCount"`
}

func NewArticle(
	authorId,
	title,
	description,
	body string,
	tagList []string,
) (*Article, error) {

	slug := slug.Make(title)

	return &Article{
		ID:             entity.NewID(),
		Slug:           slug,
		Title:          title,
		Description:    description,
		Body:           body,
		TagList:        tagList,
		Favorited:      false,
		FavoritesCount: 0,
		AuthorID:       authorId,
		CreatedAt:      time.Now().Format(time.RFC3339),
		UpdatedAt:      time.Now().Format(time.RFC3339),
	}, nil
}
