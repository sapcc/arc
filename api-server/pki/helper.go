package pki

import (
	"database/sql"
	"os"
	"path"

	"github.com/pborman/uuid"
)

// PathTo generates a path to the pki-package
func PathTo(p string) string {
	pwd, _ := os.Getwd()
	return path.Join(pwd, p)
}

// CreateTestToken save a test token in the db
func CreateTestToken(db *sql.DB, subject string) string {
	token := uuid.New()
	_, err := db.Exec("INSERT INTO tokens (id, profile, subject) VALUES($1, $2, $3)", token, "default", subject)
	if err != nil {
		panic(err)
	}
	return token
}
