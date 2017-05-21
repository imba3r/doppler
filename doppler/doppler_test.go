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
	if !s.FoundDuplicates() {
		t.Error("expected duplicates")
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
	if !s.FoundDuplicates() {
		t.Error("expected duplicates")
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
	if !s.FoundDuplicates() {
		t.Error("expected duplicates")
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

func TestMultipleTrees(t *testing.T) {
	s := doppler.NewScanner(false, false, false)
	s.ScanDirs([]string{"../tests/fixtures/tree1", "../tests/fixtures/tree3"})
	if s.HasErrors() {
		t.Error("expected no errors")
	}
	if !s.FoundDuplicates() {
		t.Error("expected duplicates")
	}
	for _, d := range s.DuplicateMap {
		if len(d.NameDuplicates) != 2 {
			t.Error("expected two name duplicate")
		}
		if len(d.HashDuplicates) != 2 {
			t.Error("expected two hash duplicate")
		}
	}
	duplicates := []string{
		"../tests/fixtures/tree3/dir1/file1.txt",
		"../tests/fixtures/tree3/dir2/file1.txt",
		"../tests/fixtures/tree1/dir2/file2.txt",
	}
	for k, v := range s.Duplicates() {
		if duplicates[k] != v {
			t.Errorf("Expected %s, got %s", duplicates[k], v)
		}
	}
}
