package entity

import (
	"github.com/sallescosta/conduit-api/pkg/entity"
	"time"
)

type Comment struct {
	ID        entity.ID `json:"id"`
	Body      string    `json:"body"`
	AuthorID  string    `json:"author_id"`
	ArticleID string    `json:"article_id"`
	CreatedAt string    `json:"created_at,omitempty"`
	UpdatedAt string    `json:"updated_at,omitempty"`
}

type AllCommentsFromAnArticle struct {
	Comments []Comment `json:"comments"`
}

func NewComment(body, authorId, articleId string) *Comment {
	return &Comment{
		ID:        entity.NewID(),
		Body:      body,
		AuthorID:  authorId,
		ArticleID: articleId,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
}
