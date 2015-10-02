// +build !integration

package swift

import (
	"flag"
	"testing"

	"github.com/codegangsta/cli"
)

//
// New()
//

func TestNewMissingParams(t *testing.T) {
	localSet := flag.NewFlagSet("test", 0)
	ctx := cli.NewContext(nil, localSet, nil)

	storage, err := New(ctx)
	if err == nil {
		t.Error("Expected to have an error")
	}
	if storage != nil {
		t.Error("Expected to have nil swift storage")
	}
}
