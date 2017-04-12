package helper

import (
	"os"

	"github.com/codingconcepts/file-wiper/model"
	"github.com/spf13/afero"
)

// GetAllFiles gets a slice of info pointers to every file in
// the collection of paths given, recursing down into them if
// they are directories
func GetAllFiles(fs afero.Fs, paths []string) (files model.Files, err error) {
	files = model.Files{}

	for _, path := range paths {
		err = afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, model.NewFile(path, info))
			}
			return nil
		})
	}

	return
}

// GetTopLevelDirectories gets a slice of info points to every directory in
// the collection of paths given
func GetTopLevelDirectories(fs afero.Fs, paths []string) (directories model.Files, err error) {
	directories = model.Files{}

	for _, path := range paths {
		var info os.FileInfo
		// break out at the first error
		if info, err = fs.Stat(path); err != nil {
			return
		}

		// skip files
		if !info.IsDir() {
			continue
		}

		directories = append(directories, model.NewFile(path, info))
	}

	return
}
