package database

import (
	"database/sql"
	"fmt"
	entityComment "github.com/sallescosta/conduit-api/internal/entity/comment"
	"log"
)

func CreateCommentsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS comments (
            id VARCHAR(255) PRIMARY KEY,
            body TEXT,
            author_id VARCHAR(255) NOT NULL,
            article_id VARCHAR(255) NOT NULL,
            createdAt TIMESTAMP DEFAULT NOW(),
            updatedAt TIMESTAMP DEFAULT NOW(),
            FOREIGN KEY (article_id) REFERENCES articles (id)
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

type CommentDB struct {
	DB *sql.DB
}

func NewComment(db *sql.DB) *CommentDB {
	return &CommentDB{DB: db}
}

func (c *CommentDB) CreateCommentDb(comment *entityComment.Comment) error {
	stmt, err := c.DB.Prepare("INSERT INTO comments (id, body, author_id, article_id, createdAt, updatedAt) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return fmt.Errorf("error preparing insert statement: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(comment.ID, comment.Body, comment.AuthorID, comment.ArticleID, comment.CreatedAt, comment.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}
