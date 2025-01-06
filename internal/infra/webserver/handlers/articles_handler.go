package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/sallescosta/conduit-api/internal/dto"
	articleEntity "github.com/sallescosta/conduit-api/internal/entity/article"
	"github.com/sallescosta/conduit-api/internal/infra/database"
	"github.com/sallescosta/conduit-api/pkg/helpers"
	"net/http"
	"strconv"
	"strings"
)

type ArticleHandler struct {
	ArticleDB database.ArticleInterface
}

func NewArticleHandler(articleDB database.ArticleInterface) *ArticleHandler {
	return &ArticleHandler{ArticleDB: articleDB}
}

func (a *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article dto.ArticleInput

	err := json.NewDecoder(r.Body).Decode(&article)
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

	art, err := articleEntity.NewArticle(
		authorId,
		article.Article.Title,
		article.Article.Description,
		article.Article.Body,
		article.Article.TagList,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating article entity\n"))
		return
	}

	err = a.ArticleDB.CreateArticle(art)
	if err != nil {
		if strings.Contains(err.Error(), "title already used") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Title already used\n"))
		} else if strings.Contains(err.Error(), "error checking title existence") {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error checking title existence\n"))
		} else if strings.Contains(err.Error(), "error preparing insert statement") {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error preparing insert statement\n"))
		} else if strings.Contains(err.Error(), "error inserting article") {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error inserting article\n"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error\n"))
		}
		return
	}

	successResponse := fmt.Sprintf("ArticleDB created successfully, title: %s", art.Title)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(successResponse))
}

func (a *ArticleHandler) ListAllArticle(w http.ResponseWriter, r *http.Request) {
	articles, err := a.ArticleDB.ListAllArticles()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := articleEntity.AllArticlesOutput{
		Articles:      articles,
		ArticlesCount: len(articles),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *ArticleHandler) FeedArticles(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")
	sort := r.URL.Query().Get("sort")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 20
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		limitInt = 0
	}

	if sort == "" {
		sort = "asc"
	}

	feed, err := a.ArticleDB.FeedArticles(limitInt, offsetInt, sort)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(feed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	article, err := a.ArticleDB.GetArticleBySlug(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(article)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (a *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var modif dto.ArticleUpdateInput

	err := json.NewDecoder(r.Body).Decode(&modif)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if modif.Article.Title == "" && modif.Article.Description == "" && modif.Article.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		//w.Write([]byte("At least one of title, description, or body must be provided"))
		//TODO: find an other way to handle errors..
		return
	}

	updatedArticle, err := a.ArticleDB.UpdateArticle(slug, modif)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(updatedArticle)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
