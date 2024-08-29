package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/sallescosta/conduit-api/internal/dto"
	userEntity "github.com/sallescosta/conduit-api/internal/entity/user"
	"github.com/sallescosta/conduit-api/internal/infra/database"
	"github.com/sallescosta/conduit-api/pkg/entity"
	"log"
	"net/http"
	"time"
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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user RegistrationInput

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	u, err := userEntity.NewUser(user.User.UserName, user.User.Email, user.User.Password)

	err = h.UserDB.CreateUser(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

type AuthenticationInput struct {
	User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user"`
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)

	var user AuthenticationInput
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("ENTRADA:", user)

	u, err := h.UserDB.FindByEmail(user.User.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Println("BYEMAIL:", u)
	if !u.ValidatePassword(user.User.Password) {
		fmt.Println("unautorized:")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	m := map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	}
	_, tokenString, _ := jwt.Encode(m)

	fmt.Println("tokenString: ", tokenString)
	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(accessToken)
	if err != nil {
		return
	}
}

func (h *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("55555555555")
	list, err := h.UserDB.GetAllUsers()
	if err != nil {
		fmt.Println("444444444444")
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
	token := r.Context()
	_, claims, _ := jwtauth.FromContext(token)
	id := claims["sub"].(string)

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
	// pegar o destinatários do request
	userName := chi.URLParam(r, "username")

	p, err := h.UserDB.GetProfileDb(userName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf(">>>>", p)

	if err := json.NewEncoder(w).Encode(p); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// pegar meu proprio id
	token := r.Context()

	_, claims, _ := jwtauth.FromContext(token)
	id := claims["sub"].(string)

	user, err := h.UserDB.FindById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	isFollowing := contain(user.Following, p.Profile.ID)

	log.Printf("1 ----", isFollowing)
	//
	//profile := dto.ProfileDTO{
	//	Profile: struct {
	//		UserName  string `json:"user_name"`
	//		Bio       string `json:"bio"`
	//		Image     string `json:"image"`
	//		Following bool   `json:"following"`
	//	}{
	//		UserName:  userName,
	//		Bio:       p.Profile.Bio,
	//		Image:     p.Profile.Image,
	//		Following: isFollowing,
	//	},
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//if err := json.NewEncoder(w).Encode(profile); err != nil {
	//	w.WriteHeader(http.StatusInternalServerError)
	//}
}

func contain(IdsSlice []entity.ID, id entity.ID) bool {
	for _, followedID := range IdsSlice {
		if followedID == id {
			return true
		}
	}
	return false
}
