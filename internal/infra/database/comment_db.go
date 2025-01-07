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
	fmt.Println("DB-authorID", comment.AuthorID)
	fmt.Println("DB-AAAArticleID", comment.ArticleID)

	_, err = stmt.Exec(comment.ID, comment.Body, comment.AuthorID, comment.ArticleID, comment.CreatedAt, comment.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (c *CommentDB) GetCommentsDb(slug string) (*entityComment.AllCommentsFromAnArticle, error) {
	articleDB := NewArticle(c.DB)
	article, err := articleDB.GetArticleBySlug(slug)

	if err != nil {
		return nil, fmt.Errorf("error getting article by slug: %w", err)
	}

	commentsQuery := `SELECT id, body, author_id, createdAt, updatedAt FROM comments WHERE article_id = $1`
	rows, err := c.DB.Query(commentsQuery, article.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying comments: %w", err)
	}
	defer rows.Close()

	var commentsList entityComment.AllCommentsFromAnArticle

	for rows.Next() {
		var comment entityComment.Comment
		err := rows.Scan(&comment.ID, &comment.Body, &comment.AuthorID, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning comment: %w", err)
		}
		commentsList.Comments = append(commentsList.Comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over comment rows: %w", err)
	}

	if len(commentsList.Comments) == 0 {
		commentsList.Comments = []entityComment.Comment{}
	}
	return &commentsList, nil
}

func (c *CommentDB) DeleteCommentsDb(id string) error {

	query := "DELETE FROM comments WHERE id = $1"
	_, err := c.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting comment: %w", err)
	}

	return nil
}
