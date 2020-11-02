package storage

import (
	"unicode"

	"github.com/fatih/structs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

func init() {
	// structs is used with squirrel (sq)
	structs.DefaultTagName = "sq"
}

type DB struct {
	*sqlx.DB
	Users UsersStorage
	Links LinksStorage
	Views ViewsStorage
}

func Open(url string) (*DB, error) {
	db, err := sqlx.Open("mysql", url)
	if err != nil {
		return nil, err
	}

	db.Mapper = reflectx.NewMapperFunc("db", toSnakeCase)

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)

	return &DB{
		DB:    db,
		Users: &Users{DB: db},
		Links: &Links{DB: db},
		Views: &Views{DB: db},
	}, nil
}

func toSnakeCase(s string) string {
	runes := []rune(s)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) {
			prev := unicode.IsLower(runes[i-1])
			next := i+1 < length && unicode.IsLower(runes[i+1])

			if prev || next {
				out = append(out, '_')
			}
		}

		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
