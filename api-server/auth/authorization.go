package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ory/ladon"
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
	RequestMethod   string
	RequestPath     string
	IdentityStatus  string
	ProjectId       string
	ProjectDomainId string
	User            User
}

type User struct {
	Id         string   `json:"id"`
	Name       string   `json:"name,omitempty"`
	DomainId   string   `json:"domain_id,omitempty"`
	DomainName string   `json:"domain_name,omitempty"`
	Roles      []string `json:"roles,omitempty"`
}

func (auth *Authorization) CheckIdentity() error {
	if auth.IdentityStatus != "Confirmed" {
		return IdentityStatusInvalid{Msg: fmt.Sprintf("%s is not 'Confirmed'", auth.IdentityStatus)}
	}

	return nil
}

// CheckPolicy checks the given roles against the policy
func (auth *Authorization) CheckPolicy(warden ladon.Ladon) error {

	for _, role := range auth.User.Roles {
		// create access request
		accessRequest := &ladon.Request{
			Subject:  role,
			Action:   strings.ToLower(auth.RequestMethod),
			Resource: auth.RequestPath,
		}

		// check request
		err := warden.IsAllowed(accessRequest)
		if err == nil {
			return nil
		}
	}

	// create error
	policies, err := warden.Manager.FindRequestCandidates(&ladon.Request{
		Action:   auth.RequestMethod,
		Resource: auth.RequestPath,
	})
	if err != nil {
		return err
	}
	roleCandidates := []string{}
	for _, pol := range policies {
		roleCandidates = append(roleCandidates, pol.GetSubjects()...)
	}

	return NotAuthorized{Msg: fmt.Sprintf("Needed roles for path %s and action %s are %s", auth.RequestPath, auth.RequestMethod, strings.Join(roleCandidates, ","))}
}

func GetIdentity(r *http.Request) *Authorization {
	return &Authorization{
		RequestMethod:   r.Method,
		RequestPath:     r.URL.Path,
		IdentityStatus:  r.Header.Get("X-Identity-Status"),
		ProjectId:       r.Header.Get("X-Project-Id"),
		ProjectDomainId: r.Header.Get("X-Project-Domain-Id"),
		User: User{
			Id:         r.Header.Get("X-User-Id"),
			Name:       r.Header.Get("X-User-Name"),
			DomainId:   r.Header.Get("X-User-Domain-Id"),
			DomainName: r.Header.Get("X-User-Domain-Name"),
			Roles:      strings.Split(r.Header.Get("X-Roles"), ","),
		},
	}
}
