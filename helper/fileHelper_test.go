package helper

import (
	"testing"

	"github.com/bassrob/file-wiper/model"
	"github.com/spf13/afero"
)

var (
	golden = map[string]string{
		"oa.txt":                  "contents of oa.txt",
		"outer_a/ia.txt":          "contents of outer_a/ia.txt",
		"outer_a/inner_a/iia.txt": "contents of outer_a/inner_a/iia.txt",

		"ob.txt":                  "contents of ob.txt",
		"outer_b/ib.txt":          "contents of outer_b/ib.txt",
		"outer_b/inner_b/iib.txt": "contents of outer_b/inner_b/iib.txt",
	}
)

func TestGetAllFiles(t *testing.T) {
	fs := setup()
	files, _ := GetAllFiles(fs, []string{"."})

	if len(files) != 6 {
		t.Fatalf("Expected 6 but got %d", len(files))
	}

	assertFile(t, files[0], "oa.txt", 18)
	assertFile(t, files[1], "ob.txt", 18)
	assertFile(t, files[2], "outer_a/ia.txt", 26)
	assertFile(t, files[3], "outer_a/inner_a/iia.txt", 35)
	assertFile(t, files[4], "outer_b/ib.txt", 26)
	assertFile(t, files[5], "outer_b/inner_b/iib.txt", 35)
}

func TestGetTopLevelDirectories(t *testing.T) {
	fs := setup()
	directories, _ := GetTopLevelDirectories(fs, []string{"outer_a", "outer_b"})

	if len(directories) != 2 {
		t.Fatalf("Expected 2 directories but got %d", len(directories))
	}

	assertDirectory(t, directories[0], "outer_a")
	assertDirectory(t, directories[1], "outer_b")
}

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
