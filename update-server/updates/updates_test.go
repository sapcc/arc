package updates

import (
	"bytes"
	"net/http"
	"testing"

	"path"
	"runtime"
	"strings"
)

func TestUpdatesNewSuccess(t *testing.T) {
	// get path to the builds
	_, filename, _, _ := runtime.Caller(0)
	buildsPath := strings.Replace(path.Dir(filename), "updates", "builds", 1)

	// get a success update
	var jsonStr = []byte(`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"arch":"amd64","os":"darwin"}}`) //
	req, _ := http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	update := New(req, buildsPath)
	if update == nil {
		t.Error("Expected update NOT to be nil. Got %q", update)
	}

	if update.Initiative != "automatically" {
		t.Error("Expected Initiative to be 'automatically'. Got %q", update.Initiative)
	}
	if update.Url != "http://0.0.0.0:3000/builds/arc_darwin_amd64_3.1.0-dev" {
		t.Error("Expected Initiative to be 'http://0.0.0.0:3000/builds/arc_darwin_amd64_3.1.0-dev'. Got %q", update.Url)
	}
	if update.Version != "3.1.0-dev" {
		t.Error("Expected Initiative to be '3.1.0-dev'. Got %q", update.Version)
	}
}

func TestUpdatesNewReturnNil(t *testing.T) {
	// post request host is missing or wrong
	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(""))
	update := New(req, "/some/build/path/")
	if update != nil {
		t.Error("Expected update to be nil. Got %q", update)
	}
	req, _ = http.NewRequest("POST", "miau", bytes.NewBufferString(""))
	update = New(req, "/some/build/path/")
	if update != nil {
		t.Error("Expected update to be nil. Got %q", update)
	}

	// post request body is empty
	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBufferString(""))
	update = New(req, "/some/build/path/")
	if update != nil {
		t.Error("Expected update to be nil. Got %q", update)
	}

	// check that the body is json
	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBufferString("not json"))
	update = New(req, "/some/build/path/")
	if update != nil {
		t.Error("Expected update to be nil. Got %q", update)
	}

	// check that the body has not the required params
	jsonStr := []byte(`{"param1":"param1"}`)
	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	update = New(req, "/some/build/path/")
	if update != nil {
		t.Error("Expected update to be nil. Got %q", update)
	}
	jsonStr = []byte(`{"app_id":"arc"}`)
	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	update = New(req, "/some/build/path/")
	if update != nil {
		t.Error("Expected update to be nil. Got %q", update)
	}
	jsonStr = []byte(`{"app_id":"arc","app_version":"0.1.0-dev"}`)
	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	update = New(req, "/some/build/path/")
	if update != nil {
		t.Error("Expected update to be nil. Got %q", update)
	}
	jsonStr = []byte(`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"os":"darwin"}}`)
	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	update = New(req, "/some/build/path/")
	if update != nil {
		t.Error("Expected update to be nil. Got %q", update)
	}

	// check wrong version format
	jsonStr = []byte(`{"app_id":"arc","app_version":"miau","tags":{"arch":"amd64","os":"darwin"}}`) //
	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	update = New(req, "/some/build/path/")
	if update != nil {
		t.Error("Expected update to be nil. Got %q", update)
	}

	// check wrong build path
	jsonStr = []byte(`{"app_id":"arc","app_version":"0.1.0-dev","tags":{"arch":"amd64","os":"darwin"}}`) //
	req, _ = http.NewRequest("POST", "http://0.0.0.0:3000/updates", bytes.NewBuffer(jsonStr))
	update = New(req, "/some/build/path/")
	if update != nil {
		t.Error("Expected update to be nil. Got %q", update)
	}
}
