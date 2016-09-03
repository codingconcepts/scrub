package helper

import (
	"os"

	"github.com/bassrob/file-wiper/model"
	"github.com/spf13/afero"
)

// GetAllFiles gets a slice of info pointers to every file in
// the collection of paths given, recursing down into them if
// they are directories
func GetAllFiles(fs afero.Fs, paths []string) (files []*model.File, err error) {
	files = []*model.File{}

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
func GetTopLevelDirectories(fs afero.Fs, paths []string) (directories []*model.File, err error) {
	directories = []*model.File{}

	for _, path := range paths {
		var info os.FileInfo
		if info, err = fs.Stat(path); err != nil || !info.IsDir() {
			return
		}

		directories = append(directories, model.NewFile(path, info))
	}

	return
}
