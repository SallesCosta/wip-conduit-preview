package database

import (
	"github.com/sallescosta/conduit-api/internal/dto"
	articleEntity "github.com/sallescosta/conduit-api/internal/entity/article"
	entityComment "github.com/sallescosta/conduit-api/internal/entity/comment"
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
	FavoriteArticleDB(slug string, isAddToFavorite bool, userID string) error
}

type ArticleInterface interface {
	CreateArticle(article *articleEntity.Article) error
	ListAllArticles() ([]articleEntity.Article, error)
	FeedArticles(limit, offset int, sort string) ([]articleEntity.Article, error)
	GetArticleBySlug(slug string) (*articleEntity.Article, error)
	UpdateArticle(slug string, article dto.ArticleUpdateInput) (*articleEntity.Article, error)
	DeleteArticleDB(slug string) error
}

type CommentInterface interface {
	CreateCommentDb(comment *entityComment.Comment) error
	GetCommentsDb(slug string) (*entityComment.AllCommentsFromAnArticle, error)
	DeleteCommentsDb(id string) error
}
