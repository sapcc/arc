package pki

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/signer"
	"github.com/pborman/uuid"
	"gitHub.***REMOVED***/monsoon/arc/api-server/auth"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"
)

type TokenRequest struct {
	signer.Subject
	Profile string
}

// CreateToken return a new sign token
func CreateToken(db *sql.DB, authorization *auth.Authorization, payload TokenRequest) (string, error) {
	// check db
	if db == nil {
		return "", errors.New("Db connection is nil")
	}

	profile := "default"
	// no need for now to change the profile
	//if payload.Profile != "" {
	//	profile = payload.Profile
	//}
	token := uuid.New()

	// At least 1 name entry and max 1 entry
	if len(payload.Subject.Names) == 0 {
		payload.Subject.Names = []csr.Name{
			csr.Name{
				OU: authorization.ProjectId,
				O:  authorization.ProjectDomainId,
			},
		}
	} else {
		// Override project and domain
		payload.Subject.Names[0].OU = authorization.ProjectId
		payload.Subject.Names[0].O = authorization.ProjectDomainId
		payload.Subject.Names[0].SerialNumber = "" // no SerialNumber in the cffsl version of arc-pki
		// just on name entry
		payload.Subject.Names = []csr.Name{payload.Subject.Names[0]}
	}

	subject, err := json.Marshal(payload.Subject)
	if err != nil {
		return "", err
	}

	// save to db
	if _, err = db.Exec(ownDb.InsertTokenQuery, token, profile, subject); err != nil {
		return "", err
	}

	return token, nil
}
