package handlers

import (
	"encoding/json"
	tagEntity "github.com/sallescosta/conduit-api/internal/entity/tag"
	"github.com/sallescosta/conduit-api/internal/infra/database"
	"net/http"
)

type TagHandler struct {
	TagDB database.TagsInterface
}

func NewTagHandler(tagDB database.TagsInterface) *TagHandler {
	return &TagHandler{TagDB: tagDB}
}

func (h *TagHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.TagDB.ListTags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if tags == nil {
		tags = []*tagEntity.Tag{}
	}
	json.NewEncoder(w).Encode(tags)
}
