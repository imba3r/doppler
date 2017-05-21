package doppler_test

import (
	"testing"

	"github.com/imba3r/doppler/doppler"
)

func TestTree1(t *testing.T) {
	s := doppler.NewScanner(false, false, false)
	s.ScanDir("../tests/fixtures/tree1")
	if s.HasErrors() {
		t.Error("expected no errors")
	}
	if len(s.DuplicateMap) != 1 {
		t.Error("expected one duplicate")
	}
	for _, d := range s.DuplicateMap {
		if len(d.NameDuplicates) > 0 {
			t.Error("expected no name duplicates")
		}
		if len(d.HashDuplicates) != 1 {
			t.Error("expected one hash duplicate")
		}
	}
}

func TestTree2(t *testing.T) {
	s := doppler.NewScanner(false, false, false)
	s.ScanDir("../tests/fixtures/tree2")
	if s.HasErrors() {
		t.Error("expected no errors")
	}
	if len(s.DuplicateMap) != 1 {
		t.Error("expected one duplicate")
	}
	for _, d := range s.DuplicateMap {
		if len(d.NameDuplicates) != 1 {
			t.Error("expected one name duplicate")
		}
		if len(d.HashDuplicates) != 1 {
			t.Error("expected one hash duplicate")
		}
	}
}

func TestTree3(t *testing.T) {
	s := doppler.NewScanner(false, false, false)
	s.ScanDir("../tests/fixtures/tree3")
	if s.HasErrors() {
		t.Error("expected no errors")
	}
	if len(s.DuplicateMap) != 1 {
		t.Error("expected one duplicate")
	}
	for _, d := range s.DuplicateMap {
		if len(d.NameDuplicates) != 1 {
			t.Error("expected one name duplicate")
		}
		if len(d.HashDuplicates) > 0 {
			t.Error("expected no hash duplicates")
		}
	}
}
