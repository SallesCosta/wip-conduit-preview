package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewComment(t *testing.T) {
	comment := NewComment(
		"body body",
		"author123",
		"article123",
	)

	assert.NotNil(t, comment)
	assert.Equal(t, "body body", comment.Body)
	assert.Equal(t, "author123", comment.AuthorID)
	assert.Equal(t, "article123", comment.ArticleID)
	assert.NotEmpty(t, comment.ID)
	assert.NotEmpty(t, comment.CreatedAt)
	assert.NotEmpty(t, comment.UpdatedAt)
}
