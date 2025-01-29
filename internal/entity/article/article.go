package entity

import (
	"github.com/gosimple/slug"
	tagEntity "github.com/sallescosta/conduit-api/internal/entity/tag"
	"github.com/sallescosta/conduit-api/pkg/entity"
	"time"
)

type Article struct {
	ID             entity.ID        `json:"id"`
	Slug           string           `json:"slug"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	Body           string           `json:"body,omitempty"`
	TagList        []*tagEntity.Tag `json:"tag_list"`
	Favorited      bool             `json:"favorited"`
	FavoritesCount uint             `json:"favorites_count"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
	AuthorID       string           `json:"author_id"`
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
	tagList []*tagEntity.Tag,
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
