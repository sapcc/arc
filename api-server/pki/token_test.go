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

var _ = Describe("Token Create", func() {

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

	It("returns a TokenBodyError error if no body is set in the request", func() {
		req, err := newAuthorizedRequest("GET", getUrl("/pki/token", url.Values{}), bytes.NewBufferString(""), map[string]string{})
		Expect(err).NotTo(HaveOccurred())
		result, err := CreateToken(db, &authorization, req)
		Expect(result).To(Equal(map[string]string{}))
		Expect(err).To(HaveOccurred())
		_, ok := err.(TokenBodyError)
		Expect(ok).To(Equal(true))
	})

	It("returns a TokenBodyError error if body is json malformated", func() {
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

	It("should let just one csr.name, all other will be removed and all SerialNumber entries set to emtpy string", func() {
		csrNames := `{"names": [{"C":"ES","OU": "projectId1", "O": "domainId1", "SerialNumber":"serialnumber1"}, {"C":"DE", "OU": "projectId2", "O": "domainId2", "SerialNumber":"serialnumber2"}] }`
		req, err := newAuthorizedRequest("GET", getUrl("/pki/token", url.Values{}), bytes.NewBufferString(csrNames), map[string]string{})
		Expect(err).NotTo(HaveOccurred())
		result, err := CreateToken(db, &authorization, req)

		var subjectData []byte
		err = db.QueryRow("SELECT subject FROM tokens WHERE id=$1", result["token"]).Scan(&subjectData)
		Expect(err).NotTo(HaveOccurred())
		var subject signer.Subject
		err = json.Unmarshal(subjectData, &subject)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(subject.Names)).To(Equal(1))
		Expect(subject.SerialNumber).To(Equal(""))                  // should be removed. Not needed in the old version of cfssl from arc-pki
		Expect(subject.Names[0].OU).To(Equal("test-project"))       // should take the OU from the authorization
		Expect(subject.Names[0].O).To(Equal("test-project-domain")) // should take the O from the authorization
		Expect(subject.Names[0].SerialNumber).To(Equal(""))         // should be removed. Not needed in the old version of cfssl from arc-pki
	})

	It("Returns a token even if empty json is send in the body request", func() {
		req, err := newAuthorizedRequest("GET", getUrl("/pki/token", url.Values{}), bytes.NewBufferString(`{}`), map[string]string{})
		Expect(err).NotTo(HaveOccurred())
		result, err := CreateToken(db, &authorization, req)

		var tokenId string
		err = db.QueryRow("SELECT id FROM tokens WHERE id=$1", result["token"]).Scan(&tokenId)
		Expect(err).NotTo(HaveOccurred())

		Expect(result["token"]).To(Equal(tokenId))
		Expect(result["url"]).To(Equal(fmt.Sprintf("http://production.***REMOVED***/api/v1/pki/sign/%s", result["token"])))
	})

	It("Returns a token", func() {
		csrNames := `{"names": [{"C":"ES","OU": "projectId1", "O": "domainId1", "SerialNumber":"serialnumber1"}, {"C":"DE", "OU": "projectId2", "O": "domainId2", "SerialNumber":"serialnumber2"}] }`
		req, err := newAuthorizedRequest("GET", getUrl("/pki/token", url.Values{}), bytes.NewBufferString(csrNames), map[string]string{})
		Expect(err).NotTo(HaveOccurred())
		result, err := CreateToken(db, &authorization, req)

		var tokenId string
		err = db.QueryRow("SELECT id FROM tokens WHERE id=$1", result["token"]).Scan(&tokenId)
		Expect(err).NotTo(HaveOccurred())

		Expect(result["token"]).To(Equal(tokenId))
		Expect(result["url"]).To(Equal(fmt.Sprintf("http://production.***REMOVED***/api/v1/pki/sign/%s", result["token"])))
	})

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
	newUrl.Host = "production.***REMOVED***"
	newUrl.Path = fmt.Sprint("/api/v1", path)
	newUrl.RawQuery = params.Encode()
	return newUrl.String()
}
