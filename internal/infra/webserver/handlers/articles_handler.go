package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/sallescosta/conduit-api/internal/dto"
	articleEntity "github.com/sallescosta/conduit-api/internal/entity/article"
	tagEntity "github.com/sallescosta/conduit-api/internal/entity/tag"
	"github.com/sallescosta/conduit-api/internal/infra/database"
	"github.com/sallescosta/conduit-api/pkg/helpers"
)

type ArticleHandler struct {
	ArticleDB database.ArticleInterface
	TagDB     database.TagsInterface
}

func NewArticleHandler(articleDB database.ArticleInterface, tagDB database.TagsInterface) *ArticleHandler {
	return &ArticleHandler{ArticleDB: articleDB, TagDB: tagDB}
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

	var lTag []*tagEntity.Tag

	for _, Itag := range article.Article.TagList {
		newTage := tagEntity.NewTag(Itag)

		lTag = append(lTag, newTage)
	}

	err = a.TagDB.CreateTag(lTag)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error creating tag\n"))
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

			return
		} else if strings.Contains(err.Error(), "error checking title existence") {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error checking title existence\n"))

			return
		} else if strings.Contains(err.Error(), "error preparing insert statement") {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error preparing insert statement\n"))

			return
		} else if strings.Contains(err.Error(), "error inserting article") {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error inserting article\n"))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error\n"))

			return
		}
		return
	}

	successResponse := fmt.Sprintf("ArticleDB created successfully, title: %s", art.Title)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(successResponse))
}

func (a *ArticleHandler) ListAllArticle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("chegou: %v", "v")))
	articles, err := a.ArticleDB.ListAllArticles()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(fmt.Appendf(nil, "Error: %v", err))
		return
	}

	if len(articles) == 0 {
		articles = []articleEntity.Article{}
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
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
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

func (a *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	err := a.ArticleDB.DeleteArticleDB(slug)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Article deleted"))
}
