package handlers

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("-->>ENTOU no ERRO<--")
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

	successResponse := fmt.Sprintf("ArticleDB created successfully, id: %s", art.ID)
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

	params := fmt.Sprintf("limit: %s,\n offset: %s", limitInt, offsetInt)

	w.Write([]byte("feedArticles.."))
	w.Write([]byte(params))
}

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
