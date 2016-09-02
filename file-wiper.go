/*
	// TODO: split code into different files and tidy up
	// TODO: progress indication
	// TODO: test with a file system mocker
*/

package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type options struct {
	sweeps           int
	debug            bool
	files            []string
	processFile      func(file string) (err error)
	processDirectory func(directory string) (err error)
}

func main() {
	opts := parse()
	processFiles(opts)

	printAndHold("Done")
}

func processFiles(opts *options) (err error) {
	for _, file := range opts.files {
		var info os.FileInfo
		if info, err = os.Stat(file); err != nil {
			return
		}

		// TODO: i don't like how nested this is
		if info.IsDir() {
			if err = processDirectory(opts, file); err != nil {
				return
			}
		} else {
			if err = processFile(opts, file); err != nil {
				return
			}
		}
	}

	return
}

func processDirectory(opts *options, directory string) (err error) {
	var paths []string
	if paths, err = getFilesRecursively(directory); err != nil {
		return
	}

	// process all of the files
	for _, path := range paths {
		if err = processFile(opts, path); err != nil {
			return
		}
	}

	// delete the directory if required
	return opts.processDirectory(directory)
}

func processFile(opts *options, file string) (err error) {
	// overwrite the file with random data
	for i := 0; i < opts.sweeps; i++ {
		if err := overwrite(file); err != nil {
			return err
		}
	}

	// delete the file if requested by the user
	return opts.processFile(file)
}

func getFilesRecursively(directory string) (files []string, err error) {
	files = []string{}

	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	return
}

func overwrite(file string) (err error) {
	var writer *bufio.Writer
	if writer, err = createWriter(file); err != nil {
		return
	}

	var info os.FileInfo
	if info, err = os.Stat(file); err != nil {
		return
	}

	// using a limited reader so we don't generate unlimited data
	limitReader := io.LimitReader(rand.Reader, info.Size())
	return pipe(limitReader, writer)
}

func createWriter(file string) (writer *bufio.Writer, err error) {
	var outputFile *os.File
	if outputFile, err = os.OpenFile(file, os.O_WRONLY, os.ModePerm); err != nil {
		return nil, err
	}

	return bufio.NewWriter(outputFile), nil
}

func pipe(reader io.Reader, writer *bufio.Writer) (err error) {
	if _, err = writer.ReadFrom(reader); err != nil {
		return
	}

	return writer.Flush()
}

func printAndHold(msg string) {
	fmt.Println(msg)
	fmt.Scan()
}

func parse() (opts *options) {
	opts = new(options)

	flag.IntVar(&opts.sweeps, "s", 10, "the number of overwrite sweeps")
	flag.BoolVar(&opts.debug, "d", false, "set flag to print files/directories, not delete them")
	flag.Parse()

	// if debug has been requested, swap in the print functionality
	if opts.debug {
		opts.processFile = debugFunc
		opts.processDirectory = debugFunc
	} else {
		opts.processFile = deleteFileFunc
		opts.processDirectory = deleteDirectoryFunc
	}

	opts.files = flag.Args()
	return
}

func deleteFileFunc(file string) (err error) {
	return os.Remove(file)
}

func deleteDirectoryFunc(directory string) (err error) {
	return os.RemoveAll(directory)
}

func debugFunc(path string) (err error) {
	fmt.Println(path)
	return nil
}
