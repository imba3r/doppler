// Package doppler scans for duplicate files by name and/or hash.
package doppler

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Scanner contains the scan options, results and occurred errors.
type Scanner struct {
	skipName bool
	skipHash bool
	absolute bool

	DuplicateMap map[string]*Duplicates `json:",omitempty"`
	ErrMap       map[string]string      `json:",omitempty"`

	nameMap map[string]string
	hashMap map[string]string
}

// NewScanner creates a new scanner with the given options.
func NewScanner(skipName bool, skipHash bool, absolute bool) *Scanner {
	s := new(Scanner)

	s.skipName = skipName
	s.skipHash = skipHash
	s.absolute = absolute

	s.nameMap = make(map[string]string)
	s.hashMap = make(map[string]string)
	s.DuplicateMap = make(map[string]*Duplicates)
	s.ErrMap = make(map[string]string)

	return s
}

// ScanDirs scans all given paths for duplicates.
func (s *Scanner) ScanDirs(paths []string) {
	for _, path := range paths {
		err := s.ScanDir(path)
		if err != nil {
			s.ErrMap[path] = err.Error()
		}
	}
}

// ScanDir scans the given path for duplicates.
func (s *Scanner) ScanDir(path string) error {
	return filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		// get absolute path if requested
		wpath := path
		if s.absolute {
			var err error
			wpath, err = filepath.Abs(path)
			if err != nil {
				s.ErrMap[path] = err.Error()
			}
		}

		// convert to sane separator
		wpath = filepath.ToSlash(wpath)

		// check for duplicate by name
		if !s.skipName {
			name := filepath.Base(wpath)
			if first, ok := s.nameMap[name]; ok {
				s.addNameDuplicate(first, wpath)
			} else {
				s.nameMap[name] = wpath
			}
		}

		// check for duplicates by content hash
		if !s.skipHash {
			hash, err := hashFile(wpath)
			if err != nil {
				s.ErrMap[path] = err.Error()
				return nil
			}
			if first, ok := s.hashMap[hash]; ok {
				s.addHashDuplicate(first, wpath)
			} else {
				s.hashMap[hash] = wpath
			}
		}
		return nil
	})
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func (s *Scanner) addHashDuplicate(path string, duplPath string) {
	if _, exists := s.DuplicateMap[path]; !exists {
		s.DuplicateMap[path] = new(Duplicates)
	}
	d := s.DuplicateMap[path]
	d.HashDuplicates = append(d.HashDuplicates, duplPath)
}

func (s *Scanner) addNameDuplicate(path string, duplPath string) {
	if _, exists := s.DuplicateMap[path]; !exists {
		s.DuplicateMap[path] = new(Duplicates)
	}
	d := s.DuplicateMap[path]
	d.NameDuplicates = append(d.NameDuplicates, duplPath)
}

// Duplicates represents the found duplicates (name/hash) of a file.
type Duplicates struct {
	NameDuplicates []string `json:",omitempty"`
	HashDuplicates []string `json:",omitempty"`
}

// HasErrors returns true if errors occurred during scanning.
func (s *Scanner) HasErrors() bool {
	return len(s.ErrMap) > 0
}

// FoundDuplicates return true if the scan found duplicates.
func (s *Scanner) FoundDuplicates() bool {
	return len(s.DuplicateMap) > 0
}

// Duplicates returns a set of all found duplicates.
func (s *Scanner) Duplicates() []string {
	set := make(map[string]struct{})
	for _, v := range s.DuplicateMap {
		var s struct{}
		for _, p := range v.NameDuplicates {
			set[p] = s
		}
		for _, p := range v.HashDuplicates {
			set[p] = s
		}
	}
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	return keys
}

// Print all found duplicates to stdout.
func (s *Scanner) Print() {
	for _, d := range s.Duplicates() {
		fmt.Println(d)
	}
}

// PrintJSON prints the scan results in JSON to stdout.
func (s *Scanner) PrintJSON() error {
	v, err := json.MarshalIndent(s.DuplicateMap, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(v))
	return nil
}
