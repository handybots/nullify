package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	UsersStorage interface {
		Create(chat Chat) error
		Exists(chat Chat) (bool, error)
		Lang(chat Chat) (string, error)
		SetLang(chat Chat, lang string) error
	}

	Users struct {
		*sqlx.DB
	}

	Chat interface {
		Recipient() string
	}

	User struct {
		CreatedAt time.Time `sq:"created_at,omitempty"`
		ID        int64     `db:"chat_id" sq:"chat_id,omitempty"`
		Lang      string    `sq:"lang,omitempty"`
		Blocked   bool      `sq:"blocked,omitempty"`
	}
)

func (db *Users) Create(chat Chat) error {
	const q = `INSERT INTO users (id, lang) VALUES (?, 'ru')`
	_, err := db.Exec(q, chat.Recipient())
	return err
}

func (db *Users) Exists(chat Chat) (has bool, _ error) {
	const q = `SELECT EXISTS(SELECT 1 FROM users WHERE id=?)`
	return has, db.Get(&has, q, chat.Recipient())
}

func (db *Users) Lang(chat Chat) (lang string, _ error) {
	const q = `SELECT lang FROM users WHERE id=?`
	return lang, db.Get(&lang, q, chat.Recipient())
}

func (db *Users) SetLang(chat Chat, lang string) error {
	const q = `UPDATE users SET lang=? WHERE id=?`
	_, err := db.Exec(q, lang, chat.Recipient())
	return err
}
