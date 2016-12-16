package pki

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/signer"
	"github.com/databus23/requestutil"
	"github.com/pborman/uuid"
	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
)

// TokenBodyError should return a http 400 error
type TokenBodyError struct {
	Msg string
}

func (e TokenBodyError) Error() string {
	return e.Msg
}

type createTokenPayload struct {
	signer.Subject
	Profile string
}

// CreateToken return a new sign token
func CreateToken(db *sql.DB, authorization *auth.Authorization, r *http.Request) (map[string]string, error) {
	// check db
	if db == nil {
		return map[string]string{}, errors.New("Db connection is nil")
	}

	// check the identity status
	err := authorization.CheckIdentity()
	if err != nil {
		return map[string]string{}, err
	}

	// read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//httpError(w, 400, err)
		return map[string]string{}, TokenBodyError{Msg: err.Error()}
	}
	r.Body.Close()

	// create payload
	var payload createTokenPayload
	if err = json.Unmarshal(body, &payload); err != nil {
		//httpError(w, 400, fmt.Errorf("Failed to parse body"))
		return map[string]string{}, TokenBodyError{Msg: "Failed to parse body"}
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
		payload.Subject.Names[0].SerialNumber = "" // no SereialNumber in the cffsl version of arc-pki
		// just on name entry
		payload.Subject.Names = []csr.Name{payload.Subject.Names[0]}
	}

	var subject []byte
	subject, err = json.Marshal(payload.Subject)
	if err != nil {
		//httpError(w, 500, err)
		return map[string]string{}, err
	}

	// save to db
	_, err = db.Exec("INSERT INTO tokens (id, profile, subject) VALUES($1, $2, $3)", token, profile, subject)
	if err != nil {
		// httpError(w, 500, err)
		return map[string]string{}, err
	}

	url := fmt.Sprintf("%s://%s/api/v1/pki/sign/%s", requestutil.Scheme(r), requestutil.HostWithPort(r), token)
	return map[string]string{"token": token, "url": url}, nil
}
