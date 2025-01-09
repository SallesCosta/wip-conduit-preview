package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"github.com/sallescosta/conduit-api/pkg/helpers"
	"log"

	"github.com/sallescosta/conduit-api/pkg/entity"

	userEntity "github.com/sallescosta/conduit-api/internal/entity/user"
)

type Error error

type ProfileWithId struct {
	Profile struct {
		ID    entity.ID
		Bio   string `json:"bio"`
		Image string `json:"image"`
	} `json:"profile"`
}

func CreateUsersTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id VARCHAR(255) PRIMARY KEY,
            username VARCHAR(50),
            email VARCHAR(50),
            password VARCHAR(255),
            bio TEXT,
            image VARCHAR(255),
            following TEXT[],
            favorites TEXT[]
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

type UserDB struct {
	DB *sql.DB
}

func NewUser(db *sql.DB) *UserDB {
	return &UserDB{DB: db}
}

func (u *UserDB) CreateUser(user *userEntity.User) error {
	stmt, err := u.DB.Prepare("INSERT INTO users (id, username, email, password, bio, image, following, " +
		"favorites) VALUES ($1," +
		" $2," +
		" $3, $4, $5, $6, $7, $8)")
	if err != nil {
		log.Fatal(err)
	}

	defer stmt.Close()

	following := make([]string, len(user.Following))
	for i, id := range user.Following {
		following[i] = id.String()
	}

	_, err = stmt.Exec(user.ID, user.UserName, user.Email, user.Password, user.Bio, user.Image, pq.Array(following),
		pq.Array(user.Favorites))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (u *UserDB) FindUserBy(field, value string) (*userEntity.User, error) {
	query := fmt.Sprintf("SELECT id, username, email, password, bio, image, following, favorites FROM users WHERE %s = $1", field)
	stmt, err := u.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var user userEntity.User

	err = stmt.QueryRow(value).Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.Bio, &user.Image,
		pq.Array(&user.Following), pq.Array(&user.Favorites))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserDB) FindByEmail(email string) (*userEntity.User, error) {
	return u.FindUserBy("email", email)
}

func (u *UserDB) FindById(id string) (*userEntity.User, error) {
	return u.FindUserBy("id", id)
}

func (u *UserDB) UpdateUserDb(email, username, password, image, bio string) (*userEntity.User, error) {
	user, err := u.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	if username != "" {
		user.UserName = username
	}

	if image != "" {
		user.Image = image
	}

	if bio != "" {
		user.Bio = bio
	}

	hashedPass, err := userEntity.DoHash(password)
	if password != "" {
		user.Password = string(hashedPass)
	}

	stmt, err := u.DB.Prepare("UPDATE users SET username = $1, password = $2, image = $3, bio = $4 WHERE email = $5")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.UserName, user.Password, user.Image, user.Bio, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserDB) GetAllUsers() ([]userEntity.User, error) {
	rows, err := u.DB.Query("SELECT id, username, email, password, bio, image, following, favorites FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []userEntity.User
	for rows.Next() {
		var user userEntity.User

		err := rows.Scan(&user.ID, &user.UserName, &user.Email, &user.Password, &user.Bio, &user.Image,
			pq.Array(&user.Following), pq.Array(&user.Favorites))
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		users = []userEntity.User{}
	}

	return users, nil
}

func (u *UserDB) GetProfileDb(userName string) (*ProfileWithId, error) {
	query := "SELECT id, bio, image FROM users WHERE username = $1"
	stmt, err := u.DB.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var profile ProfileWithId

	err = stmt.QueryRow(userName).Scan(&profile.Profile.ID, &profile.Profile.Bio, &profile.Profile.Image)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("profile not found: %w", err)
		}
		return nil, err
	}

	return &profile, nil
}

func (u *UserDB) UpdateFollowingUserDb(id string, following []entity.ID) error {
	stmt, err := u.DB.Prepare("UPDATE users SET following = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	followingStr := make([]string, len(following))
	for i, id := range following {
		followingStr[i] = id.String()
	}

	_, err = stmt.Exec(pq.Array(followingStr), id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserDB) FavoriteArticleDB(slug string, isAddToFavorite bool, userID string) error {
	articleDB := NewArticle(u.DB)
	articleToFavorite, err := articleDB.GetArticleBySlug(slug)
	if err != nil {
		return err
	}

	articleId := articleToFavorite.ID

	user, err := u.FindById(userID)
	if err != nil {
		return err
	}

	if isAddToFavorite {
		for _, id := range user.Favorites {
			if id == articleId {
				return fmt.Errorf("article already in favorites")
			}
		}
		user.Favorites = append(user.Favorites, articleId)

		articleToFavorite.FavoritesCount++
	} else {

		originalLength := len(user.Favorites)

		user.Favorites = helpers.RemoveItem(user.Favorites, articleId)
		if len(user.Favorites) == originalLength {
			return fmt.Errorf("article not found in favorites")
		}

		if articleToFavorite.FavoritesCount > 0 {
			articleToFavorite.FavoritesCount--
		} else {
			return fmt.Errorf("favoritesCount is zero")
		}
	}

	stmt, err := u.DB.Prepare("UPDATE users SET favorites = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	favoritesStr := make([]string, len(user.Favorites))
	for i, id := range user.Favorites {
		favoritesStr[i] = id.String()
	}

	_, err = stmt.Exec(pq.Array(favoritesStr), userID)
	if err != nil {
		return err
	}

	stmt, err = u.DB.Prepare("UPDATE articles SET favoritesCount = $1 WHERE id = $2")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(articleToFavorite.FavoritesCount, articleId.String())
	if err != nil {
		return err
	}

	return nil
}
