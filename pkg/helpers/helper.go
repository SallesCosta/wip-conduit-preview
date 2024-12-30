package helpers

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	"github.com/sallescosta/conduit-api/pkg/entity"
	"net/http"
)

func Contain(IdsSlice []entity.ID, id entity.ID) bool {
	for _, followedID := range IdsSlice {
		if followedID == id {
			return true
		}
	}
	return false
}

func GetMyOwnIdbyToken(r *http.Request) (string, error) {
	token := r.Context()

	_, claims, err := jwtauth.FromContext(token)
	if err != nil {
		return "", err
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("token claims do not contain 'sub'")
	}

	return sub, nil
}
