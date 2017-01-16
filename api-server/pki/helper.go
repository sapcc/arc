package pki

import (
	"database/sql"

	"github.com/cloudflare/cfssl/csr"
	"github.com/pborman/uuid"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
)

// CreateTestToken save a test token in the db
func CreateTestToken(db *sql.DB, subject string) string {
	token := uuid.New()
	_, err := db.Exec(ownDb.InsertTokenQuery, token, "default", subject)
	if err != nil {
		panic(err)
	}
	return token
}

// GetTestToken get a saved token
func GetTestToken(db *sql.DB, token string) (string, error) {
	var profile string
	var subjectData []byte
	err := db.QueryRow(ownDb.GetTokenQuery, token).Scan(&profile, &subjectData)
	if err != nil {
		return "", err
	}
	return profile, nil
}

// CreateCSR creates a privateKey and PEM encoded csr
func CreateCSR(commonName, organization, organizationalUnit string) (csreq, key []byte, err error) {

	req := csr.New()
	req.CN = commonName
	req.Names = []csr.Name{csr.Name{O: organization, OU: organizationalUnit}}
	csreq, key, err = csr.ParseRequest(req)
	return
}
