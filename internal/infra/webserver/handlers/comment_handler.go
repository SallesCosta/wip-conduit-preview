package handlers

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("authorID", authorId)
	fmt.Println("AAAArticleID", articleId)
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

	commentsList, err := c.CommentDB.GetCommentsDb(slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(commentsList)
}

func (c *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := c.CommentDB.DeleteCommentsDb(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("comment removed."))
}
