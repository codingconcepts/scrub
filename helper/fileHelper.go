package helper

import (
	"os"
	"path/filepath"

	"github.com/bassrob/file-wiper/model"
)

// GetAllFiles gets a slice of info pointers to every file in
// the collection of paths given, recursing down into them if
// they are directories
func GetAllFiles(paths []string) (files []*model.File, err error) {
	files = []*model.File{}

	for _, path := range paths {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
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
func GetTopLevelDirectories(paths []string) (directories []*model.File, err error) {
	directories = []*model.File{}

	for _, path := range paths {
		var info os.FileInfo
		if info, err = os.Stat(path); err != nil || !info.IsDir() {
			return
		}

		directories = append(directories, model.NewFile(path, info))
	}

	return
}
