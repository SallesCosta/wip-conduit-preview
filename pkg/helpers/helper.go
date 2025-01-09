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

func RemoveItem[T comparable](s []T, x T) []T {
	result := make([]T, 0, len(s))
	for _, item := range s {
		if item != x {
			result = append(result, item)
		}
	}
	return result
}

func Ternary[T any](val bool, x, y T) T {
	if val {
		return x
	}
	return y
}
