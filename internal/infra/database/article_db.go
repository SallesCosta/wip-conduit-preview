package database

import (
	"database/sql"
	"fmt"
	slugMaker "github.com/gosimple/slug"
	"github.com/lib/pq"
	"github.com/sallescosta/conduit-api/internal/dto"
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
	var exists bool
	err := a.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM articles WHERE title = $1)", article.Title).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking title existence: %w", err)
	}
	if exists {
		return fmt.Errorf("title already used")
	}

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
	rows, err := a.DB.Query("SELECT id, author_id, slug, title, description, body, favorited, favoritesCount, tag_list, createdAt, updatedAt FROM articles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []articleEntity.Article

	for rows.Next() {
		var article articleEntity.Article

		err := rows.Scan(&article.ID, &article.AuthorID, &article.Slug, &article.Title, &article.Description,
			&article.Body, &article.Favorited, &article.FavoritesCount, pq.Array(&article.TagList), &article.CreatedAt, &article.UpdatedAt)

		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

func (a *ArticleDB) FeedArticles(limit, offset int, sort string) ([]articleEntity.Article, error) {
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	query := fmt.Sprintf(`SELECT id, author_id, slug, title, description, body, favorited, favoritesCount, tag_list, createdAt, updatedAt FROM articles ORDER BY createdAt %s LIMIT $1 OFFSET $2`, sort)

	rows, err := a.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedArticles []articleEntity.Article

	for rows.Next() {
		var article articleEntity.Article

		err := rows.Scan(&article.ID, &article.AuthorID, &article.Slug, &article.Title, &article.Description,
			&article.Body, &article.Favorited, &article.FavoritesCount, pq.Array(&article.TagList), &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			fmt.Println("erro no scan", err)
			return nil, err
		}

		feedArticles = append(feedArticles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return feedArticles, nil
}

func (a *ArticleDB) GetArticleBySlug(slug string) (*articleEntity.Article, error) {
	query := "SELECT id, author_id, slug, title, description, body, favorited, favoritesCount, tag_list, createdAt, updatedAt FROM articles WHERE slug = $1"
	stmt, err := a.DB.Prepare(query)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	article := &articleEntity.Article{}

	err = stmt.QueryRow(slug).Scan(&article.ID, &article.AuthorID, &article.Slug, &article.Title, &article.Description,
		&article.Body, &article.Favorited, &article.FavoritesCount, pq.Array(&article.TagList), &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (a *ArticleDB) UpdateArticle(slug string, article dto.ArticleUpdateInput) (*articleEntity.Article, error) {
	articleToUpdate, err := a.GetArticleBySlug(slug)

	if err != nil {
		return nil, err
	}

	if article.Article.Title != "" {
		articleToUpdate.Title = article.Article.Title
		articleToUpdate.Slug = slugMaker.Make(article.Article.Title)
	}

	if article.Article.Description != "" {
		articleToUpdate.Description = article.Article.Description
	}

	if article.Article.Body != "" {
		articleToUpdate.Body = article.Article.Body
	}

	stmt, err := a.DB.Prepare("UPDATE articles SET title = $1, description = $2, body = $3,  favorited = $4 WHERE slug = $5")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(articleToUpdate.Title, articleToUpdate.Description, articleToUpdate.Body, articleToUpdate.Favorited, slug)
	if err != nil {
		return nil, err
	}

	return articleToUpdate, nil
}

func (a *ArticleDB) DeleteArticleDB(slug string) error {
	article, err := a.GetArticleBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("article not found")
		}
		return fmt.Errorf("error checking article existence: %w", err)
	}

	query := "DELETE FROM articles WHERE slug = $1"
	_, err = a.DB.Exec(query, article.Slug)
	if err != nil {
		return fmt.Errorf("error deleting article: %w", err)
	}
	return nil
}

//func (a *ArticleDB) FavoriteArticleDB(slug string, isAddToFavorite bool) error {
//	articleToFavorite, err := a.GetArticleBySlug(slug)
//	if err != nil {
//		fmt.Println("-->erro aqui<--", err)
//		return err
//	}
//
//	if articleToFavorite.Favorited == isAddToFavorite {
//		fmt.Println("-->erro A<--", err)
//		return fmt.Errorf("article already in the desired favorite state")
//	}
//
//	articleToFavorite.Favorited = isAddToFavorite
//
//	if isAddToFavorite {
//		articleToFavorite.FavoritesCount++
//	} else {
//		if articleToFavorite.FavoritesCount > 0 {
//			articleToFavorite.FavoritesCount--
//		} else {
//
//			fmt.Println("-->erro B<--")
//			return fmt.Errorf("favoritesCount cannot be negative")
//		}
//	}
//
//	stmt, err := a.DB.Prepare("UPDATE articles SET favorited = $1, favoritesCount = $2 WHERE slug = $3")
//	if err != nil {
//		fmt.Println("erro 2 --->", err)
//		return err
//	}
//	defer stmt.Close()
//
//	_, err = stmt.Exec(articleToFavorite.Favorited, articleToFavorite.FavoritesCount, slug)
//	if err != nil {
//		fmt.Println("erro 3 --->", err)
//		return err
//	}
//
//	return nil
//}
