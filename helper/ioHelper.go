package helper

import (
	"bufio"
	"io"
	"os"

	"github.com/bassrob/file-wiper/model"
)

// CreateWriter creates a *bufio writer from a given file info pointer,
// opening the file with WRONLY permissions, ready for writing to
func CreateWriter(file *model.File) (writer *bufio.Writer, err error) {
	var outputFile *os.File
	if outputFile, err = os.OpenFile(file.FullPath, os.O_WRONLY, os.ModePerm); err != nil {
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
