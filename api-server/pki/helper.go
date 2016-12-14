package pki

import (
	"database/sql"
	"os"
	"path"

	"github.com/pborman/uuid"
)

func PathTo(p string) string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, p)
}

func CreateTestToken(db *sql.DB, subject string) string {
	token := uuid.New()
	_, err := db.Exec("INSERT INTO tokens (id, profile, subject) VALUES($1, $2, $3)", token, "default", subject)
	if err != nil {
		panic(err)
	}
	return token
}
