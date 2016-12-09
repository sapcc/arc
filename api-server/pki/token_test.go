// +build integration

package pki_test

import (
	auth "gitHub.***REMOVED***/monsoon/arc/api-server/authorization"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/pki"

	"github.com/cloudflare/cfssl/signer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var _ = Describe("Tokens", func() {

	var (
		authorization = auth.Authorization{}
	)

	JustBeforeEach(func() {
		// reset authorization
		authorization.IdentityStatus = "Confirmed"
		authorization.User = auth.User{Id: "userID", Name: "Arturo", DomainId: "monsoon2_id", DomainName: "monsoon_name"}
		authorization.ProjectId = "test-project"
		authorization.ProjectDomainId = "test-project-domain"
	})

	It("returns an error if no db connection is given", func() {
		req, err := newAuthorizedRequest("GET", getUrl("/pki/token", url.Values{}), bytes.NewBufferString("{}"), map[string]string{})
		Expect(err).NotTo(HaveOccurred())
		result, err := CreateToken(nil, &authorization, req)
		Expect(err).To(HaveOccurred())
		Expect(result).To(Equal(map[string]string{}))
	})

	It("returns an error if no body is set in the request", func() {
		req, err := newAuthorizedRequest("GET", getUrl("/pki/token", url.Values{}), bytes.NewBufferString(""), map[string]string{})
		Expect(err).NotTo(HaveOccurred())
		result, err := CreateToken(db, &authorization, req)
		Expect(result).To(Equal(map[string]string{}))
		Expect(err).To(HaveOccurred())
		_, ok := err.(TokenBodyError)
		Expect(ok).To(Equal(true))
	})

	It("returns an error if body is json malformated", func() {
		req, err := newAuthorizedRequest("GET", getUrl("/pki/token", url.Values{}), bytes.NewBufferString(`{"CN": "agent name"`), map[string]string{})
		Expect(err).NotTo(HaveOccurred())
		result, err := CreateToken(db, &authorization, req)
		Expect(result).To(Equal(map[string]string{}))
		Expect(err).To(HaveOccurred())
		_, ok := err.(TokenBodyError)
		Expect(ok).To(Equal(true))
	})

	It("should set a csr.name with project id and domain id if json request body empty", func() {
		req, err := newAuthorizedRequest("GET", getUrl("/pki/token", url.Values{}), bytes.NewBufferString(`{}`), map[string]string{})
		Expect(err).NotTo(HaveOccurred())
		result, err := CreateToken(db, &authorization, req)

		var profile string
		var subjectData []byte
		err = db.QueryRow("SELECT profile, subject FROM tokens WHERE id=$1", result["token"]).Scan(&profile, &subjectData)
		Expect(err).NotTo(HaveOccurred())
		var subject signer.Subject
		err = json.Unmarshal(subjectData, &subject)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(subject.Names)).To(Equal(1))
		Expect(subject.Names[0].OU).To(Equal("test-project"))
		Expect(subject.Names[0].O).To(Equal("test-project-domain"))
	})

	//It("Returns a token", func() {})

})

func newAuthorizedRequest(method, urlStr string, body io.Reader, headerOptions map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Identity-Status", `Confirmed`)
	req.Header.Add("X-Project-Id", `test-project`)
	req.Header.Add("X-User-Id", `arc_test`)
	req.Header.Add("X-Project-Domain-Id", `test-project-domain`)

	// add extra headers
	for k, v := range headerOptions {
		req.Header.Add(k, v)
	}

	return req, nil
}

func getUrl(path string, params url.Values) string {
	//var newUrl *url.URL
	newUrl, _ := url.Parse("")
	newUrl.Path = fmt.Sprint("/api/v1", path)
	newUrl.RawQuery = params.Encode()
	return newUrl.String()
}
