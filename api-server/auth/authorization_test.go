package auth

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/ory/ladon"
	ladon_mem "github.com/ory/ladon/manager/memory"
)

func TestCheckPolicyActionMatchSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://arc.app/api/v1/jobs", strings.NewReader(""))
	req.Header.Set("X-Roles", "automation_viewer")

	warden := &ladon.Ladon{Manager: ladon_mem.NewMemoryManager()}
	warden.Manager.Create(&ladon.DefaultPolicy{
		ID:          "2",
		Description: "viewers can call get in any resource",
		Subjects:    []string{"automation_viewer"},
		Actions:     []string{"get"},
		Resources:   []string{"<.*>"},
		Effect:      ladon.AllowAccess,
	})

	authorization := GetIdentity(req)
	err := authorization.CheckPolicy(*warden)
	if err != nil {
		t.Error(fmt.Sprintf("Expected to not have an error, but got %s", err.Error()))
	}
}

func TestCheckPolicyResourceMatchSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://arc.app/api/v1/jobs/1234566789/logs", strings.NewReader(""))
	req.Header.Set("X-Roles", "automation_viewer")

	warden := &ladon.Ladon{Manager: ladon_mem.NewMemoryManager()}
	warden.Manager.Create(&ladon.DefaultPolicy{
		ID:          "2",
		Description: "viewers can call get in any resource",
		Subjects:    []string{"automation_viewer"},
		Actions:     []string{"get"},
		Resources:   []string{`<^(\/api\/v[0-9]\/)(jobs\/)(.*)(\/logs)$>`},
		Effect:      ladon.AllowAccess,
	})

	authorization := GetIdentity(req)
	err := authorization.CheckPolicy(*warden)
	if err != nil {
		t.Error(fmt.Sprintf("Expected to not have an error, but got %s", err.Error()))
	}
}

func TestCheckPolicyNoRolesGiven(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://arc.app/api/v1/jobs/1234566789/logs", strings.NewReader(""))
	req.Header.Set("X-Roles", "")

	warden := &ladon.Ladon{Manager: ladon_mem.NewMemoryManager()}
	warden.Manager.Create(&ladon.DefaultPolicy{
		ID:          "1",
		Description: "viewers can call get in any resource",
		Subjects:    []string{"automation_viewer"},
		Actions:     []string{"get"},
		Resources:   []string{`<.*>`},
		Effect:      ladon.AllowAccess,
	})

	authorization := GetIdentity(req)
	err := authorization.CheckPolicy(*warden)
	if err == nil {
		t.Error(fmt.Sprintf("Expected to not have an error, but got %s", err.Error()))
	}
}

func TestCheckPolicyFailingRole(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://arc.app/api/v1/jobs/1234566789/logs", strings.NewReader(""))
	req.Header.Set("X-Roles", "keystone_admin,compute_viewer")

	warden := &ladon.Ladon{Manager: ladon_mem.NewMemoryManager()}
	warden.Manager.Create(&ladon.DefaultPolicy{
		ID:          "1",
		Description: "viewers can call get in any resource",
		Subjects:    []string{"automation_admin,god"},
		Actions:     []string{"<.*>"},
		Resources:   []string{`<.*>`},
		Effect:      ladon.AllowAccess,
	})
	warden.Manager.Create(&ladon.DefaultPolicy{
		ID:          "2",
		Description: "viewers can call get in any resource",
		Subjects:    []string{"automation_viewer"},
		Actions:     []string{"get"},
		Resources:   []string{`<.*>`},
		Effect:      ladon.AllowAccess,
	})

	authorization := GetIdentity(req)
	err := authorization.CheckPolicy(*warden)
	if err == nil {
		t.Error(fmt.Sprintf("Expected to have an error, but got %s", err.Error()))
	}
}
