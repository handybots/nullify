package storage

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	zeroWidthRune0 = '\u200c'
	zeroWidthRune1 = '\u200d'
)

type (
	LinksStorage interface {
		Create(link Link) (Link, error)
		ByID(id int64) (Link, error)
		ByString(s string) (Link, error)
		ByUserID(chat Chat) ([]Link, error)
		SetDeleted(link Link) error
		CountByUserID(chat Chat) (int, error)
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

// String returns generated zero-width characters chain by the ID seed.
func (l Link) String() string {
	var b strings.Builder
	b.Grow(64)
	for i := 0; i != 64; i++ {
		if l.ID&(1<<i) != 0 {
			b.WriteRune(zeroWidthRune0)
		} else {
			b.WriteRune(zeroWidthRune1)
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

func (db *Links) ByID(id int64) (link Link, _ error) {
	const q = `SELECT * FROM links WHERE id = ?`
	return link, db.Get(&link, q, id)
}

func (db *Links) ByString(s string) (link Link, _ error) {
	var id, b int64 = 0, 1
	for _, r := range []rune(s) {
		if r == zeroWidthRune0 {
			id += b
		}
		b = b << 1
	}

	return db.ByID(id)
}

func (db *Links) ByUserID(chat Chat) (links []Link, _ error) {
	const q = `SELECT * FROM links WHERE user_id = ? AND deleted = false`
	return links, db.Select(&links, q, chat.Recipient())
}

func (db *Links) SetDeleted(link Link) error {
	const q = "UPDATE links SET deleted = ? WHERE id = ?"
	_, err := db.Exec(q, link.Deleted, link.ID)
	return err
}

func (db *Links) CountByUserID(chat Chat) (number int, _ error) {
	const q = `SELECT COUNT(1) FROM links WHERE user_id = ? AND deleted = false`
	return number, db.Get(&number, q, chat.Recipient())
}
