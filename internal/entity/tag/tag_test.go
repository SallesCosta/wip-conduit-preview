package entity

import (
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestNewTag(t *testing.T) {
  tag := NewTag("Go")
  
  assert.NotNil(t, tag)
  assert.NotNil(t, tag.ID)
  assert.Equal(t, "Go", tag.Name)
}
