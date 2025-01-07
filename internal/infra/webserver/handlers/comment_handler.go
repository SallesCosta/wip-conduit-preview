package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/sallescosta/conduit-api/internal/dto"
	entityComment "github.com/sallescosta/conduit-api/internal/entity/comment"
	"github.com/sallescosta/conduit-api/internal/infra/database"
	"github.com/sallescosta/conduit-api/pkg/helpers"
	"net/http"
)

type CommentHandler struct {
	CommentDB database.CommentInterface
}

func NewCommentHandler(commentDB database.CommentInterface) *CommentHandler {
	return &CommentHandler{CommentDB: commentDB}
}

func (c *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment dto.AddCommentInput

	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorId, err := helpers.GetMyOwnIdbyToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	articleId := comment.Comment.ArticleID

	newComment := entityComment.NewComment(
		comment.Comment.Body,
		authorId,
		articleId,
	)

	err = c.CommentDB.CreateCommentDb(newComment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating comment entity\n"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newComment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (c *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	commentsList, err := c.CommentDB.GetComments(slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(commentsList)
}
