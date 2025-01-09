package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/sallescosta/conduit-api/internal/dto"
	userEntity "github.com/sallescosta/conduit-api/internal/entity/user"
	"github.com/sallescosta/conduit-api/internal/infra/database"
	"github.com/sallescosta/conduit-api/pkg/helpers"
)

type UserHandler struct {
	UserDB database.UserInterface
}

func NewUserHandler(userDB database.UserInterface) *UserHandler {
	return &UserHandler{UserDB: userDB}
}

type RegistrationInput struct {
	User struct {
		UserName string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user"`
}

type AuthenticationInput struct {
	User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user RegistrationInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	email := user.User.Email
	userFound, err := h.UserDB.FindByEmail(email)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if userFound != nil {
		http.Error(w, "UserDB already exists", http.StatusConflict)
		return
	}

	u, err := userEntity.NewUser(user.User.UserName, user.User.Email, user.User.Password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = h.UserDB.CreateUser(u)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)

	var user AuthenticationInput
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u, err := h.UserDB.FindByEmail(user.User.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	if !u.ValidatePassword(user.User.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	m := map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	}
	_, tokenString, _ := jwt.Encode(m)

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(accessToken)

	if err != nil {
		return
	}
}

func (h *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	list, err := h.UserDB.GetAllUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.UserDTO
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updatedUser, err := h.UserDB.UpdateUserDb(user.User.Email, user.User.UserName, user.User.Password, user.User.Image, user.User.Bio)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(updatedUser)
	if err != nil {
		return
	}
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.GetMyOwnIdbyToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	user, err := h.UserDB.FindById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) GetProfileUser(w http.ResponseWriter, r *http.Request) {
	userName := chi.URLParam(r, "username")

	if userName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := h.UserDB.GetProfileDb(userName)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "UserDB not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	id, err := helpers.GetMyOwnIdbyToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	user, err := h.UserDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	isFollowing := helpers.Contain(user.Following, p.Profile.ID)

	profile := dto.ProfileDTO{
		Profile: dto.Profile{
			UserName:  userName,
			Bio:       p.Profile.Bio,
			Image:     p.Profile.Image,
			Following: isFollowing,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	myId, err := helpers.GetMyOwnIdbyToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	mySelf, err := h.UserDB.FindById(myId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	userName := chi.URLParam(r, "username")

	p, err := h.UserDB.GetProfileDb(userName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	if mySelf.ID == p.Profile.ID {
		w.WriteHeader(http.StatusConflict)
		return
	}

	var isFollowing bool

	switch r.Method {
	case http.MethodPost:

		if helpers.Contain(mySelf.Following, p.Profile.ID) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		mySelf.Following = append(mySelf.Following, p.Profile.ID)
		isFollowing = true
	case http.MethodDelete:
		for i, followedID := range mySelf.Following {
			if followedID == p.Profile.ID {
				mySelf.Following = append(mySelf.Following[:i], mySelf.Following[i+1:]...)
				break
			}
		}
		isFollowing = false
	}

	err = h.UserDB.UpdateFollowingUserDb(mySelf.ID.String(), mySelf.Following)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	profile := dto.ProfileDTO{
		Profile: dto.Profile{
			UserName:  userName,
			Bio:       p.Profile.Bio,
			Image:     p.Profile.Image,
			Following: isFollowing,
		},
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *UserHandler) FavoriteArticle(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	isAddToFavorite := r.Method == http.MethodPost
 
	id, err := helpers.GetMyOwnIdbyToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	err = h.UserDB.FavoriteArticleDB(slug, isAddToFavorite, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("--->>Error: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	addMessage := "Article added to favorites"
	removedMessage := "Article removed from favorites"

	var message string
	if isAddToFavorite {
		message = addMessage
	} else {
		message = removedMessage
	}
	w.Write([]byte(message))
}
