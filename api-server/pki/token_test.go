// +build integration

package pki_test

import (
	"time"

	auth "gitHub.***REMOVED***/monsoon/arc/api-server/auth"
	. "gitHub.***REMOVED***/monsoon/arc/api-server/pki"

	"github.com/cloudflare/cfssl/signer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
	ownDb "gitHub.***REMOVED***/monsoon/arc/api-server/db"

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
		var tr TokenRequest
		err := json.Unmarshal([]byte("{}"), &tr)
		Expect(err).NotTo(HaveOccurred())

		result, err := CreateToken(nil, &authorization, tr)
		Expect(err).To(HaveOccurred())
		Expect(result).To(Equal(""))
	})

	It("should set a csr.name with project id and domain id if json request body empty", func() {
		var tr TokenRequest
		err := json.Unmarshal([]byte(`{}`), &tr)
		Expect(err).NotTo(HaveOccurred())
		token, err := CreateToken(db, &authorization, tr)
		Expect(err).NotTo(HaveOccurred())

		var profile string
		var subjectData []byte
		err = db.QueryRow("SELECT profile, subject FROM tokens WHERE id=$1", token).Scan(&profile, &subjectData)
		Expect(err).NotTo(HaveOccurred())
		var subject signer.Subject
		err = json.Unmarshal(subjectData, &subject)
		Expect(err).NotTo(HaveOccurred())

		Expect(len(subject.Names)).To(Equal(1))
		Expect(subject.Names[0].OU).To(Equal("test-project"))
		Expect(subject.Names[0].O).To(Equal("test-project-domain"))
	})

	It("should let just one csr.name all others will be removed, override O and OU and all SerialNumber entries set to emtpy string", func() {
		csrNames := `{"names": [{"C":"ES","OU": "projectId1", "O": "domainId1", "SerialNumber":"serialnumber1"}, {"C":"DE", "OU": "projectId2", "O": "domainId2", "SerialNumber":"serialnumber2"}] }`
		var tr TokenRequest
		err := json.Unmarshal([]byte(csrNames), &tr)
		Expect(err).NotTo(HaveOccurred())
		token, err := CreateToken(db, &authorization, tr)
		Expect(err).NotTo(HaveOccurred())

		var subjectData []byte
		err = db.QueryRow("SELECT subject FROM tokens WHERE id=$1", token).Scan(&subjectData)
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

	It("should set the common name provided", func() {
		var tr TokenRequest
		err := json.Unmarshal([]byte(`{"CN": "agent-test-name"}`), &tr)
		Expect(err).NotTo(HaveOccurred())
		token, err := CreateToken(db, &authorization, tr)
		Expect(err).NotTo(HaveOccurred())

		var subjectData []byte
		err = db.QueryRow("SELECT subject FROM tokens WHERE id=$1", token).Scan(&subjectData)
		Expect(err).NotTo(HaveOccurred())
		var subject signer.Subject
		err = json.Unmarshal(subjectData, &subject)
		Expect(err).NotTo(HaveOccurred())
		Expect(subject.CN).To(Equal("agent-test-name"))
	})

	It("refuses to create tokens with funny common names", func() {
		var tr TokenRequest
		err := json.Unmarshal([]byte(`{"CN": "agent-test-name; DROP TABLE users;"}`), &tr)
		Expect(err).NotTo(HaveOccurred())
		_, err = CreateToken(db, &authorization, tr)
		Expect(err).To(Equal(ErrorInvalidCommonName))
	})

	It("Returns a token even if empty json is send in the body request", func() {
		var tr TokenRequest
		err := json.Unmarshal([]byte(`{}`), &tr)
		Expect(err).NotTo(HaveOccurred())
		token, err := CreateToken(db, &authorization, tr)
		Expect(err).NotTo(HaveOccurred())

		var tokenId string
		err = db.QueryRow("SELECT id FROM tokens WHERE id=$1", token).Scan(&tokenId)
		Expect(err).NotTo(HaveOccurred())
		Expect(token).To(Equal(tokenId))
	})

	It("should create a token even if the token request is empty", func() {
		tr := TokenRequest{}
		_, err := CreateToken(db, &authorization, tr)
		Expect(err).NotTo(HaveOccurred())
	})

	It("Returns a token", func() {
		csrNames := `{"names": [{"C":"ES","OU": "projectId1", "O": "domainId1", "SerialNumber":"serialnumber1"}, {"C":"DE", "OU": "projectId2", "O": "domainId2", "SerialNumber":"serialnumber2"}] }`
		var tr TokenRequest
		err := json.Unmarshal([]byte(csrNames), &tr)
		Expect(err).NotTo(HaveOccurred())
		token, err := CreateToken(db, &authorization, tr)
		Expect(err).NotTo(HaveOccurred())

		var tokenId string
		err = db.QueryRow("SELECT id FROM tokens WHERE id=$1", token).Scan(&tokenId)
		Expect(err).NotTo(HaveOccurred())

		Expect(token).To(Equal(tokenId))
	})

})

var _ = Describe("PruneTokens", func() {

	It("returns an error if no db connection is given", func() {
		occurrencies, err := PruneTokens(nil)
		Expect(err).To(HaveOccurred())
		Expect(occurrencies).To(Equal(int64(0)))
	})

	It("should clean pki tokens older than 1 hour", func() {
		// insert old token
		_, err := db.Exec(ownDb.InsertTokenWithCreatedAtQuery, uuid.New(), "default", `{}`, time.Now().Add((-65)*time.Minute))
		Expect(err).NotTo(HaveOccurred())

		occurrencies, err := PruneTokens(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(occurrencies).To(Equal(int64(1)))
	})

	It("should NOT clean pki tokens older than 1 hour", func() {
		// insert old token
		_, err := db.Exec(ownDb.InsertTokenWithCreatedAtQuery, uuid.New(), "default", `{}`, time.Now().Add((-15)*time.Minute))
		Expect(err).NotTo(HaveOccurred())

		occurrencies, err := PruneTokens(db)
		Expect(err).NotTo(HaveOccurred())
		Expect(occurrencies).To(Equal(int64(0)))
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
	newUrl.Host = "production.app"
	newUrl.Path = fmt.Sprint("/api/v1", path)
	newUrl.RawQuery = params.Encode()
	return newUrl.String()
}
