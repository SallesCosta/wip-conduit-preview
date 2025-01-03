package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/sallescosta/conduit-api/internal/infra/database"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/sallescosta/conduit-api/cmd/server/router"
	"github.com/sallescosta/conduit-api/configs"
)

func main() {
	config, err := configs.LoadConfig()

	if err != nil {
		panic(err)
	}

	connStr := fmt.Sprintf(
		"user=%s password= %s dbname=%s sslmode=disable",
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		slog.Error("Error opening database connection", slog.String("error", err.Error()))
	}

	err = db.Ping()
	if err != nil {
		slog.Info("Verify the docker. Open the docker and run `docker-compose up -d`")
	} else {
		fmt.Println("Successfully connected to the database!")
		database.CreateUsersTable(db)
		database.CreateArticlesTable(db)
	}

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()

	router.Init(r, config, db)
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
