package helper

import (
	"os"
	"testing"

	"github.com/bassrob/file-wiper/model"
	"github.com/spf13/afero"
)

func TestOverwriteModifiesFile(t *testing.T) {
	fs := setup()
	files, _ := GetAllFiles(fs, []string{"."})

	for _, file := range files {
		Overwrite(fs, file)

		content, _ := afero.ReadFile(fs, file.FullPath)
		actual := string(content)
		expected := golden[file.FullPath]

		if actual == expected {
			t.Fatal("File was not overwritten")
		}
	}
}

func TestProcessFileDoesNotDeleteIfDebug(t *testing.T) {
	fs := setup()
	opts := model.NewOptions(fs, 1, true)
	files, _ := GetAllFiles(fs, []string{"."})

	for _, file := range files {
		ProcessFile(fs, opts, file)

		content, _ := afero.ReadFile(fs, file.FullPath)
		actual := string(content)
		expected := golden[file.FullPath]

		if actual == expected {
			t.Fatal("File was not overwritten")
		}
	}
}

func TestProcessFileDeletesIfNotDebug(t *testing.T) {
	fs := setup()
	opts := model.NewOptions(fs, 1, false)
	files, _ := GetAllFiles(fs, []string{"."})

	for _, file := range files {
		ProcessFile(fs, opts, file)

		if _, err := afero.ReadFile(fs, file.FullPath); !os.IsNotExist(err) {
			t.Fatal("File was not overwritten")
		}
	}
}
