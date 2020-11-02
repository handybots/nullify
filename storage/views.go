package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	ViewsStorage interface {
		Create(view View) error
	}

	Views struct {
		*sqlx.DB
	}

	View struct {
		ID        int64     `db:"id" sq:"chat_id,omitempty"`
		CreatedAt time.Time `sq:"created_at"`
		LinkID    int64     `sq:"link_id"`
		IP        string    `sq:"ip"`
		UserAgent string    `sq:"user_agent"`
	}
)

func (db *Views) Create(view View) error {
	const q = `INSERT INTO views (link_id, ip, user_agent) VALUES (?, ?, ?)`
	_, err := db.Exec(q, view.LinkID, view.IP, view.UserAgent)
	return err
}
