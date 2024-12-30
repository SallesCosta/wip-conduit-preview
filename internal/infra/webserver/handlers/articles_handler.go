package handlers

import (
	"encoding/json"
	"github.com/sallescosta/conduit-api/internal/dto"
	articleEntity "github.com/sallescosta/conduit-api/internal/entity/article"
	"github.com/sallescosta/conduit-api/internal/infra/database"
	"github.com/sallescosta/conduit-api/pkg/helpers"
	"net/http"
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

	err = a.ArticleDB.CreateArticle(art)
	if err != nil {
		if strings.Contains(err.Error(), "title already used") {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("Title already used"))
		}
		if strings.Contains(err.Error(), "error checking title existence") {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error checking title existence"))

		}
		if strings.Contains(err.Error(), "error preparing insert statement") {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error preparing insert statement"))
		}

		if strings.Contains(err.Error(), "error inserting article") {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error inserting article"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Article created successfully"))
}

//func GetArticles(w http.ResponseWriter, r *http.Request) {
//	tag := r.URL.Query().Get("tag")
//	author := r.URL.Query().Get("author")
//	favorited := r.URL.Query().Get("favorited")
//	limit := r.URL.Query().Get("limit")
//	offset := r.URL.Query().Get("offset")
//
//	params := fmt.Sprintf("tag: %s,\n author: %s,\n favorited: %s,\n limit: %s,\n offset: %s", tag, author, favorited,
//		limit, offset)
//
//	w.Write([]byte("getArticles"))
//	w.Write([]byte(params))
//}
//
//func GetArticlesFeed(w http.ResponseWriter, r *http.Request) {
//	url := r.URL.String()
//	limit := r.URL.Query().Get("limit")
//	offset := r.URL.Query().Get("offset")
//
//	params := fmt.Sprintf("limit: %s,\n offset: %s", limit, offset)
//	response := fmt.Sprintf("url: %s - GetArticlesFeed", url)
//
//	w.Write([]byte(response))
//	w.Write([]byte(params))
//}
