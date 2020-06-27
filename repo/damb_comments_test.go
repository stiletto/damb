package repo

import (
	"testing"

	"github.com/sebdah/goldie/v2"
)

func TestDambComments(t *testing.T) {
	g := goldie.New(t)
	opts, err := LoadDambOptions("testdata/Dockerfile")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(opts)
	g.Assert(t, "TestDambComments.Context", []byte(opts.Context))
}
