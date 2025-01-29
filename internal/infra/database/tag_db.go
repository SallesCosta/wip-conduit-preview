package database

import (
	"database/sql"
	"fmt"
	tagEntity "github.com/sallescosta/conduit-api/internal/entity/tag"
	"log"
)

func CreateTagsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS Tags (
            id VARCHAR(255) PRIMARY KEY,
            name VARCHAR(255) UNIQUE NOT NULL
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Tags table created")
	return nil
}

func CreateArticleTagsTable(db *sql.DB) error {
	query := `
        CREATE TABLE IF NOT EXISTS ArticleTags (
            article_id VARCHAR(255) NOT NULL,
            tag_id VARCHAR(255) NOT NULL,
            PRIMARY KEY (article_id, tag_id),
            FOREIGN KEY (article_id) REFERENCES Articles (id) ON DELETE CASCADE,
            FOREIGN KEY (tag_id) REFERENCES Tags (id) ON DELETE CASCADE
        );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("ArticleTags table created")
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
		stmt, err := t.DB.Prepare("INSERT INTO Tags (id, name) VALUES ($1, $2)")
		if err != nil {
			fmt.Println("error preparing insert statement")
			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(tag.ID, tag.Name)
		if err != nil {
			fmt.Println("error preparing insert statement")
			return err
		}
	}

	return nil
}

func (t *TagDB) ListTags() ([]*tagEntity.Tag, error) {
	rows, err := t.DB.Query("SELECT name FROM Tags")
	if err != nil {
		return nil, err
	}

	var tags []*tagEntity.Tag

	for rows.Next() {
		var tag *tagEntity.Tag
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
