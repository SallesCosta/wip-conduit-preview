package database

import (
	"database/sql"
	"fmt"
	"log"

	articleEntity "github.com/sallescosta/conduit-api/internal/entity/article"
)

func CreateArticlesTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS articles (
            id VARCHAR(255) PRIMARY KEY,
            slug VARCHAR(100) UNIQUE NOT NULL,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            body TEXT,
            createdAt TIMESTAMP DEFAULT NOW(),
            updatedAt TIMESTAMP DEFAULT NOW(),
            favorited BOOLEAN DEFAULT FALSE,
            favoritesCount INT DEFAULT 0,
            author_id VARCHAR(255) NOT NULL,
            FOREIGN KEY (author_id) REFERENCES users (id)
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

type Article struct {
	DB *sql.DB
}

func NewArticle(db *sql.DB) *Article {
	return &Article{DB: db}
}

func (a *Article) CreateArticle(article *articleEntity.Article) error {
	var existingTitle string

	err := a.DB.QueryRow("SELECT title FROM articles WHERE title = $1", article.Title).Scan(&existingTitle)
	if err == nil {
		return fmt.Errorf("title already used: %s", existingTitle)
	}
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error checking title existence: %w", err)
	}

	stmt, err := a.DB.Prepare(`
        INSERT INTO articles ( id, author_id, slug, title, description, body, createdAt, updatedAt )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
    `)
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		article.ID,
		article.AuthorID,
		article.Slug,
		article.Title,
		article.Description,
		article.Body,
		article.CreatedAt,
		article.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error inserting article: %w", err)
	}
	return nil
}

func (a *Article) GetAllArticles() ([]articleEntity.Article, error) {
	rows, err := a.DB.Query("SELECT id, author_id, slug, title, description, body, createdAt, updatedAt FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []articleEntity.Article
 
	for rows.Next() {
		var article articleEntity.Article
		err := rows.Scan(&article.ID, &article.AuthorID, &article.Slug, &article.Title, &article.Description, &article.Body, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}
  
		articles = append(articles, article)
	}
	return articles, nil
}
