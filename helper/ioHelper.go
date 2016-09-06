package helper

import (
	"bufio"
	"crypto/rand"
	"io"
	"os"

	"github.com/bassrob/file-wiper/model"
	"github.com/spf13/afero"
)

// ProcessFile overwrites a given file with random data a given
// number of times
func ProcessFile(fs afero.Fs, opts *model.Options, file *model.File) (err error) {
	// overwrite the file with random data
	for i := 0; i < opts.Sweeps; i++ {
		if err := Overwrite(fs, file); err != nil {
			return err
		}
	}

	// delete the file if requested by the user
	return opts.ProcessFile(file)
}

// Overwrites a given file with random data
func Overwrite(fs afero.Fs, file *model.File) (err error) {
	var writer *bufio.Writer
	if writer, err = CreateWriter(fs, file); err != nil {
		return
	}

	// using a limited reader so we don't generate unlimited data
	limitReader := io.LimitReader(rand.Reader, file.Size)
	return Pipe(limitReader, writer)
}

// CreateWriter creates a *bufio writer from a given file info pointer,
// opening the file with WRONLY permissions, ready for writing to
func CreateWriter(fs afero.Fs, file *model.File) (writer *bufio.Writer, err error) {
	var outputFile afero.File
	if outputFile, err = fs.OpenFile(file.FullPath, os.O_WRONLY, os.ModePerm); err != nil {
		return nil, err
	}

	return bufio.NewWriter(outputFile), nil
}

// Pipe copies the data from an io reader into a *bufio writer,
// flusing when done
func Pipe(reader io.Reader, writer *bufio.Writer) (err error) {
	if _, err = writer.ReadFrom(reader); err != nil {
		return
	}

	return writer.Flush()
}
