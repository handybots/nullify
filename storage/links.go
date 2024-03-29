package storage

import (
	"time"

	"github.com/handybots/nullify/linkgen"
	"github.com/jmoiron/sqlx"
)

type (
	LinksStorage interface {
		Create(link Link) (Link, error)
		ByID(id int64) (Link, error)
		ByString(s string) (Link, error)
		SetDeleted(id int64, deleted bool) error
		ByUser(user Chat) ([]Link, error)
		CountByUser(user Chat) (int, error)
		ExistsForUser(id int64, user Chat) (bool, error)
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

func (l Link) String() string {
	return linkgen.S(l.ID)
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
	const q = `SELECT * FROM links WHERE id = ? AND deleted = false`
	return link, db.Get(&link, q, id)
}

func (db *Links) ByString(s string) (link Link, _ error) {
	return db.ByID(linkgen.N(s))
}

func (db *Links) SetDeleted(id int64, deleted bool) error {
	const q = "UPDATE links SET deleted = ? WHERE id = ?"
	_, err := db.Exec(q, deleted, id)
	return err
}

func (db *Links) ByUser(user Chat) (links []Link, _ error) {
	const q = `SELECT * FROM links WHERE user_id = ? AND deleted = false`
	return links, db.Select(&links, q, user.Recipient())
}

func (db *Links) CountByUser(user Chat) (number int, _ error) {
	const q = `SELECT COUNT(*) FROM links WHERE user_id = ? AND deleted = false`
	return number, db.Get(&number, q, user.Recipient())
}

func (db *Links) ExistsForUser(id int64, user Chat) (has bool, _ error) {
	const q = `SELECT EXISTS(SELECT 1 FROM links WHERE id=? AND user_id=? AND deleted=false)`
	return has, db.Get(&has, q, id, user.Recipient())
}
