package model

import (
	"fmt"
	"os"
)

// File contains basic information about a file or a directory,
// storing the values, rather than exposing them through methods
// as os.FileInfo does
type File struct {
	FullPath string
	Size     int64
	IsDir    bool
}

func (file *File) String() string {
	return fmt.Sprintf("%s (%v): %d", file.FullPath, file.IsDir, file.Size)
}

// NewFile creates a pointer to a File struct from an os.FileInfo
func NewFile(path string, info os.FileInfo) (file *File) {
	return &File{
		FullPath: path,
		Size:     info.Size(),
		IsDir:    info.IsDir(),
	}
}

// TotalSize sums the sizes of a given set of File structs
func TotalSize(files []*File) (size int64) {
	for _, file := range files {
		size += file.Size
	}
	return
}
