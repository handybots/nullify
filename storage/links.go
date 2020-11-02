package storage

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	LinksStorage interface {
		Create(link Link) (Link, error)
		ByUserID(chat Chat) ([]Link, error)
	}

	Links struct {
		*sqlx.DB
	}

	Link struct {
		CreatedAt time.Time `sq:"created_at,omitempty"`
		ID        int64     `db:"id" sq:"id,omitempty"`
		UserID    int       `sq:"user_id,omitempty"`
		URL       string    `sq:"url,omitempty"`
		Deleted   bool      `sq:"deleted,omitempty"`
	}
)

func LinkFromString(s string) int64 {
	var n, b int64 = 0, 1
	for _, r := range []rune(s) {
		if r == '\u200c' {
			n += b
		}
		b = b << 1
	}
	return n
}

// String returns generated zero-width characters chain by the ID seed.
func (l Link) String() string {
	var b strings.Builder
	b.Grow(64)
	for i := 0; i != 64; i++ {
		if l.ID&(1<<i) != 0 {
			b.WriteString("\u200c")
		} else {
			b.WriteString("\u200d")
		}
	}
	return b.String()
}

func (db *Links) Create(link Link) (Link, error) {
	const q = `INSERT INTO links (user_id, url) VALUES (?, ?)`

	r, err := db.Exec(q, link.UserID, link.URL)
	if err != nil {
		return link, err
	}

	link.ID, err = r.LastInsertId()
	return link, err
}

func (db *Links) ByUserID(chat Chat) (links []Link, _ error) {
	const q = `SELECT * FROM links WHERE user_id = ? LIMIT 5`
	return links, db.Select(&links, q, chat.Recipient())
}
