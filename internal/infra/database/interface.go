package database

import (
	articleEntity "github.com/sallescosta/conduit-api/internal/entity/article"
	userEntity "github.com/sallescosta/conduit-api/internal/entity/user"
	"github.com/sallescosta/conduit-api/pkg/entity"
)

type UserInterface interface {
	CreateUser(user *userEntity.User) error
	FindByEmail(email string) (*userEntity.User, error)
	FindById(id string) (*userEntity.User, error)
	GetAllUsers() ([]userEntity.User, error)
	UpdateUserDb(email, username, password, image, bio string) (*userEntity.User, error)
	GetProfileDb(userName string) (*ProfileWithId, error)
	UpdateFollowingUserDb(id string, following []entity.ID) error
}

type ArticleInterface interface {
	CreateArticle(article *articleEntity.Article) error
	ListAllArticles() ([]articleEntity.Article, error)
	FeedArticles(limit, offset int, sort string) ([]articleEntity.Article, error)
}
