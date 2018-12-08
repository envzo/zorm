package example

import (
	"testing"

	"github.com/envzo/zorm/gen"
)

func TestGen(t *testing.T) {
	if err := gen.Gen("../example/pod_user.yaml", "../example", "orm"); err != nil {
		t.Fatal(err)
	}
}
