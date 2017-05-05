package helper

import (
	"path/filepath"
	"testing"

	"github.com/codingconcepts/scrub/model"
	"github.com/spf13/afero"
)

var (
	golden = map[string]string{
		"oa.txt": "contents of oa.txt",
		filepath.Join("outer_a", "ia.txt"):             filepath.Join("contents of outer_a", "ia.txt"),
		filepath.Join("outer_a", "inner_a", "iia.txt"): filepath.Join("contents of outer_a", "inner_a", "iia.txt"),

		"ob.txt": "contents of ob.txt",
		filepath.Join("outer_b", "ib.txt"):             filepath.Join("contents of outer_b", "ib.txt"),
		filepath.Join("outer_b", "inner_b", "iib.txt"): filepath.Join("contents of outer_b", "inner_b", "iib.txt"),
	}
)

func setup() (fs afero.Fs) {
	fs = afero.NewMemMapFs()

	for key, value := range golden {
		file, _ := fs.Create(key)
		file.WriteString(value)
	}

	return
}

func assertFile(t *testing.T, file *model.File, expectedFullPath string, expectedSize int64) {
	if file.FullPath != expectedFullPath {
		t.Fatalf("Expected %s but got %s", expectedFullPath, file.FullPath)
	}
	if file.Size != expectedSize {
		t.Fatalf("Expected %d but got %d", expectedSize, file.Size)
	}
	if file.IsDir {
		t.Fatalf("Expected a file but got a directory")
	}
}

func assertDirectory(t *testing.T, file *model.File, expectedFullPath string) {
	if file.FullPath != expectedFullPath {
		t.Fatalf("Expected %s but got %s", expectedFullPath, file.FullPath)
	}
	if !file.IsDir {
		t.Fatalf("Expected a directory but got a file")
	}
}
