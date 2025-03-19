package router

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/sallescosta/conduit-api/configs"
	"github.com/sallescosta/conduit-api/internal/infra/database"

	"github.com/sallescosta/conduit-api/internal/infra/webserver/handlers"
)

func Init(r *chi.Mux, config *configs.Conf, db *sql.DB) {
	userHandler := handlers.NewUserHandler(database.NewUser(db))
	articleHandler := handlers.NewArticleHandler(database.NewArticle(db), database.NewTag(db))
	commentHandler := handlers.NewCommentHandler(database.NewComment(db))
	tagHandler := handlers.NewTagHandler(database.NewTag(db))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", config.JwtExpiresIn))

	r.Post("/api/users", userHandler.CreateUser)
	r.Post("/api/users/login", userHandler.GetJWT)

	r.Route("/api/user", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Put("/", userHandler.UpdateUser)
		r.Get("/", userHandler.GetCurrentUser)
		r.Get("/all", userHandler.ListAllUsers)
	})

	r.Route("/api/profiles", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/{username}", userHandler.GetProfileUser)
		r.Post("/{username}/follow", userHandler.FollowUser)
		r.Delete("/{username}/follow", userHandler.FollowUser)
	})

	r.Route("/api/articles", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Post("/", articleHandler.CreateArticle)
		r.Get("/", articleHandler.ListAllArticle)
		r.Get("/feed", articleHandler.FeedArticles)

		r.Get("/{slug}", articleHandler.GetArticle)
		r.Put("/{slug}", articleHandler.UpdateArticle)
		r.Delete("/{slug}", articleHandler.DeleteArticle)

		r.Post("/{slug}/favorite", userHandler.FavoriteArticle)
		r.Delete("/{slug}/favorite", userHandler.FavoriteArticle)

		r.Post("/{slug}/comments", commentHandler.CreateComment)
		r.Get("/{slug}/comments", commentHandler.GetComments)
		r.Delete("/{slug}/comments/{id}", commentHandler.DeleteComment)
	})

	r.Route("/api/tags", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/", tagHandler.ListTags)
	})
}
