package authorization

import (
	"fmt"
	"net/http"
)

var IdentityStatusInvalid = fmt.Errorf("Authorization: Identity status invalid.")
var NotAuthorized = fmt.Errorf("Authorization: not authorized.")

type Authorization struct {
	IdentityStatus string
	UserId         string
	ProjectId      string
}

func (auth *Authorization) CheckIdentity() error {
	if auth.IdentityStatus != "Confirmed" {
		return IdentityStatusInvalid
	}

	return nil
}

func GetIdentity(r *http.Request) *Authorization {
	return &Authorization{
		IdentityStatus: r.Header.Get("X-Identity-Status"),
		UserId:         r.Header.Get("X-User-Id"),
		ProjectId:      r.Header.Get("X-Project-Id"),
	}
}
