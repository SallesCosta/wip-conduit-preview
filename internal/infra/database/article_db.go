package database

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"

	articleEntity "github.com/sallescosta/conduit-api/internal/entity/article"
)

func CreateArticlesTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS articles (
            id VARCHAR(255) PRIMARY KEY,
            author_id VARCHAR(255) NOT NULL,
            slug VARCHAR(100) UNIQUE NOT NULL,
            title VARCHAR(255) NOT NULL,
            description TEXT,
            body TEXT,
            tag_list TEXT[],
            favorited BOOLEAN DEFAULT FALSE,
            favoritesCount INT DEFAULT 0,
            createdAt TIMESTAMP DEFAULT NOW(),
            updatedAt TIMESTAMP DEFAULT NOW(),
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

type ArticleDB struct {
	DB *sql.DB
}

func NewArticle(db *sql.DB) *ArticleDB {
	return &ArticleDB{DB: db}
}

func (a *ArticleDB) CreateArticle(article *articleEntity.Article) error {
	fmt.Println("tagList", article.TagList)

	stmt, err := a.DB.Prepare(`
		INSERT INTO articles (
			id, author_id, slug, title, description, body, favorited, favoritesCount, tag_list, createdAt, updatedAt
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
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
		article.Favorited,
		article.FavoritesCount,
		pq.Array(article.TagList),
		article.CreatedAt,
		article.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error inserting article: %w", err)
	}
	return nil
}

func (a *ArticleDB) ListAllArticles() ([]articleEntity.Article, error) {
	rows, err := a.DB.Query("SELECT id, author_id, slug, title, description, body, tag_list, createdAt, " +
		"updatedAt FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []articleEntity.Article

	for rows.Next() {
		var article articleEntity.Article
		err := rows.Scan(&article.ID, &article.AuthorID, &article.Slug, &article.Title, &article.Description,
			&article.Body, &article.TagList, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}
	return articles, nil
}

func (a *ArticleDB) FeedArticles(limit, offset int, sort string) ([]articleEntity.Article, error) {
	if sort != "asc" || sort != "esc" {
		sort = "asc"
	}

	query := `SELECT id, slug, title, description, body, favorited, favorited_count, tag_list, created_at, updated_at,
 author_id FROM articles`

	rows, err := a.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}

	var feedArticles []articleEntity.Article
	for rows.Next() {
		var article articleEntity.Article
		rows.Scan(&article.ID, &article.Slug, &article.Title, &article.Description, &article.Body, &article.Favorited,
			&article.FavoritesCount, &article.TagList, &article.CreatedAt, &article.UpdatedAt)
		feedArticles = append(feedArticles, article)
	}

	return feedArticles, nil

}
