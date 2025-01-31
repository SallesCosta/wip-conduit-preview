package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tags = []string{"Go", "Tailwind", "Templ", "HTMX"}

func TestNewArticle(t *testing.T) {
	article, err := NewArticle(
		"author123",
		"My title",
		"My description",
		"Body, body, body...",
		tags)

	assert.Nil(t, err)
	assert.NotNil(t, article)
	assert.NotEmpty(t, article.ID)
	assert.Equal(t, "author123", article.AuthorID)
	assert.Equal(t, "My title", article.Title)
	assert.Equal(t, "My description", article.Description)
	assert.Equal(t, "Body, body, body...", article.Body)
	assert.Equal(t, []string{"Go", "Tailwind", "Templ", "HTMX"}, article.TagList)
}
