package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	LinksStorage interface {
		Create(link Link) error
		ByUserID(id int) ([]Link, error)
	}

	Links struct {
		*sqlx.DB
	}

	Link struct {
		ID               int64     `db:"id" sq:"chat_id,omitempty"`
		UserID           int       `sq:"user_id"`
		CreatedAt        time.Time `sq:"created_at,omitempty"`
		SourceLink       string    `sq:"source_link"`
		Link             string    `sq:"link"`
		Conversion       int       `sq:"conversion"`
		ConversionUnique int       `sq:"conversion_unique"`
	}
)

func (db *Links) Create(link Link) error {
	const q = `INSERT INTO links (user_id, source_link, link) VALUES (?, ?, ?)`
	_, err := db.Exec(q, link.UserID, link.SourceLink, link.Link)
	return err
}

func (db *Links) ByUserID(id int) ([]Link, error) {
	var links []Link
	const q = `SELECT * FROM links WHERE user_id = ? LIMIT 5`
	return links, db.Select(&links, q, id)
}
