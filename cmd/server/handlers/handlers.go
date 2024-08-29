package handlers

import (
	"fmt"
	"io"
	"net/http"
)

func GenericHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	urlFromRquest := r.URL.String()
	if err != nil {
		panic("eitaNois")
	}

	w.Write([]byte(urlFromRquest))
	w.Write(body)
}

func GetArticles(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	author := r.URL.Query().Get("author")
	favorited := r.URL.Query().Get("favorited")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	params := fmt.Sprintf("tag: %s,\n author: %s,\n favorited: %s,\n limit: %s,\n offset: %s", tag, author, favorited,
		limit, offset)

	w.Write([]byte("getArticles"))
	w.Write([]byte(params))
}

func GetArticlesFeed(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	params := fmt.Sprintf("limit: %s,\n offset: %s", limit, offset)
	response := fmt.Sprintf("url: %s - GetArticlesFeed", url)

	w.Write([]byte(response))
	w.Write([]byte(params))
}
