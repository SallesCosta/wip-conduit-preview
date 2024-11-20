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
		//r.Delete("/follow", handlers.GenericHandler)
	})

	r.Route("/api/articles", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)

		//r.Get("/", handlers.GetArticles)
		//r.Post("/", handlers.GenericHandler)
		//r.Get("/feed", handlers.GetArticlesFeed)
		//r.Get("/{slug}", handlers.GenericHandler)
		//r.Put("/{slug}", handlers.GenericHandler)
		//r.Delete("/{slug}", handlers.GenericHandler)

		//r.Post("/{slug}/comments", handlers.GenericHandler)
		//r.Get("/{slug}/comments", handlers.GenericHandler) // authentication optional, return multiple comments
		//r.Delete("/{slug}/comments/{id}", handlers.GenericHandler)
	})

	//r.Get("/api/tags", handlers.GenericHandler)
}
