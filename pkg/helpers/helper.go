package helpers

import (
	"github.com/go-chi/jwtauth"
	"github.com/sallescosta/conduit-api/pkg/entity"
	"net/http"
)

//interface Helpers = {
//
//}

func Contain(IdsSlice []entity.ID, id entity.ID) bool {
	for _, followedID := range IdsSlice {
		if followedID == id {
			return true
		}
	}
	return false
}

func GetMyOwnIdbyToken(r *http.Request) string {
	token := r.Context()

	_, claims, _ := jwtauth.FromContext(token)
	return claims["sub"].(string)
}
