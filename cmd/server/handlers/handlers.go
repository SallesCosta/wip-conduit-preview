package handlers

import (
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

