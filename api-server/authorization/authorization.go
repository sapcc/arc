package authorization

import (
	"fmt"
	"net/http"
)

// var IdentityStatusInvalid = fmt.Errorf("Authorization: Identity status invalid.")
// var NotAuthorized = fmt.Errorf("Authorization: not authorized.")

type IdentityStatusInvalid struct {
	Msg string
}

func (e IdentityStatusInvalid) Error() string {
	return fmt.Sprint("Authorization: Identity status invalid. ", e.Msg)
}

type NotAuthorized struct {
	Msg string
}

func (e NotAuthorized) Error() string {
	return fmt.Sprint("Authorization: not authorized. ", e.Msg)
}

type Authorization struct {
	IdentityStatus string
	ProjectId      string
	User           User
}

type User struct {
	Id         string `json:"id"`
	Name       string `json:"name,omitempty"`
	DomainId   string `json:"domain_id,omitempty"`
	DomainName string `json:"domain_name,omitempty"`
}

func (auth *Authorization) CheckIdentity() error {
	if auth.IdentityStatus != "Confirmed" {
		return IdentityStatusInvalid{Msg: fmt.Sprintf("%s is not 'Confirmed'", auth.IdentityStatus)}
	}

	return nil
}

func GetIdentity(r *http.Request) *Authorization {
	return &Authorization{
		IdentityStatus: r.Header.Get("X-Identity-Status"),
		ProjectId:      r.Header.Get("X-Project-Id"),
		User: User{
			Id:         r.Header.Get("X-User-Id"),
			Name:       r.Header.Get("X-User-Name"),
			DomainId:   r.Header.Get("X-User-Domain-Id"),
			DomainName: r.Header.Get("X-User-Domain-Name"),
		},
	}
}
