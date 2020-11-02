package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	ViewsStorage interface {
		Create(view View) error
		Stats(linkID int64) (ViewStats, error)
	}

	Views struct {
		*sqlx.DB
	}

	View struct {
		CreatedAt time.Time `sq:"created_at,omitempty"`
		ID        int64     `db:"id" sq:"chat_id,omitempty"`
		LinkID    int64     `sq:"link_id,omitempty"`
		IP        string    `sq:"ip,omitempty"`
		UserAgent string    `sq:"user_agent,omitempty"`
	}

	ViewStats struct {
		Views  int `db:"views"`
		Unique int `db:"unique"`
	}
)

func (db *Views) Create(view View) error {
	const q = `INSERT INTO views (link_id, ip, user_agent) VALUES (?, ?, ?)`
	_, err := db.Exec(q, view.LinkID, view.IP, view.UserAgent)
	return err
}

func (db *Views) Stats(linkID int64) (stats ViewStats, _ error) {
	const q = `
		SELECT
			COUNT(*) "views",
			COUNT(DISTINCT(ip)) "unique"
		FROM views
		WHERE link_id=?`

	return stats, db.Get(&stats, q, linkID)
}
