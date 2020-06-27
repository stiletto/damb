package dambfile

import (
	"testing"

	"github.com/sebdah/goldie/v2"
)

func TestLoader(t *testing.T) {
	cfg, err := FindAndLoad("testdata")
	if err != nil {
		t.Fatal(err)
	}
	err = cfg.Recompute()
	if err != nil {
		t.Fatal("Recompute", err)
	}
	t.Log(cfg)
	g := goldie.New(t)
	g.AssertJson(t, "TestLoader.cfg", cfg)
	// t.Fatal("not implemented")
}
