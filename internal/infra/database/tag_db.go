package database

import (
	"database/sql"
	"fmt"
	tagEntity "github.com/sallescosta/conduit-api/internal/entity/tag"
	"log"
)

func CreateTagsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS tags (
            id VARCHAR(255) PRIMARY KEY,
            name VARCHAR(255) UNIQUE NOT NULL
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("tags table created")
	return nil
}

func CreateArticleTagsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS article_tags (
            article_id VARCHAR(255) NOT NULL,
            tag_id VARCHAR(255) NOT NULL,
            PRIMARY KEY (article_id, tag_id),
            FOREIGN KEY (article_id) REFERENCES articles (id) ON DELETE CASCADE,
            FOREIGN KEY (tag_id) REFERENCES tags (id) ON DELETE CASCADE
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("article_tags table created")
	return nil
}

type TagDB struct {
	DB *sql.DB
}

func NewTag(db *sql.DB) *TagDB {
	return &TagDB{DB: db}
}

func (t *TagDB) CreateTag(tags []*tagEntity.Tag) error {
	for _, tag := range tags {
		var exists bool
		_ = t.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tags WHERE name = $1)", tag.Name).Scan(&exists)
		if exists {
			continue
		}

		stmt, err := t.DB.Prepare("INSERT INTO tags (id, name) VALUES ($1, $2)")
		if err != nil {
			fmt.Println("error preparing insert statement")
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(tag.ID, tag.Name)
		if err != nil {
			fmt.Println("error inserting tag")
			return err
		}
	}

	return nil
}

func (t *TagDB) ListTags() ([]*tagEntity.Tag, error) {
	rows, err := t.DB.Query("SELECT id, name FROM tags")
	if err != nil {
		fmt.Errorf("error querying tags: %v", err)
		return nil, err
	}

	defer rows.Close()

	var tags []*tagEntity.Tag

	for rows.Next() {
		tag := &tagEntity.Tag{}

		err = rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			fmt.Errorf("error scan", err)
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
