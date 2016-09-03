/*
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

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/bassrob/file-wiper/helper"
	"github.com/bassrob/file-wiper/model"
)

func main() {
	opts := parse()

	if err := process(opts); err != nil {
		fmt.Printf("An error occurred: %s", err)
	}

	printAndHold("Finished, press any key to quit")
}

func process(opts *model.Options) (err error) {
	var files []*model.File
	if files, err = helper.GetAllFiles(opts.Files); err != nil {
		return
	}

	var directories []*model.File
	if directories, err = helper.GetTopLevelDirectories(opts.Files); err != nil {
		return
	}

	totalSize := model.TotalSize(files)
	progressBar := createProgressBar(totalSize)

	// loop through all of the nested files and process them
	progressBar.Start()
	for _, file := range files {
		if err := processFile(opts, file); err != nil {
			panic(err)
		}
		progressBar.Add64(file.Size)
	}
	progressBar.Finish()

	// now that all of the nested files have been cleaned up,
	// remove the top-level directories
	for _, file := range directories {
		if err := opts.ProcessDirectory(file); err != nil {
			panic(err)
		}
	}

	return
}

func createProgressBar(totalSize int64) (bar *pb.ProgressBar) {
	bar = pb.New64(totalSize)
	bar.SetUnits(pb.U_BYTES)

	return
}

func processFile(opts *model.Options, file *model.File) (err error) {
	// overwrite the file with random data
	for i := 0; i < opts.Sweeps; i++ {
		if err := overwrite(file); err != nil {
			return err
		}
	}

	// delete the file if requested by the user
	return opts.ProcessFile(file)
}

func overwrite(file *model.File) (err error) {
	var writer *bufio.Writer
	if writer, err = helper.CreateWriter(file); err != nil {
		return
	}

	// using a limited reader so we don't generate unlimited data
	limitReader := io.LimitReader(rand.Reader, file.Size)
	return helper.Pipe(limitReader, writer)
}

func printAndHold(msg string) {
	fmt.Println(msg)
	fmt.Scanln()
}

func parse() (opts *model.Options) {
	opts = new(model.Options)

	flag.IntVar(&opts.Sweeps, "s", 10, "the number of overwrite sweeps")
	flag.BoolVar(&opts.Debug, "d", false, "set flag to print files/directories, not delete them")
	flag.Parse()

	// if debug has been requested, swap in the print functionality
	if opts.Debug {
		opts.ProcessFile = debugFunc
		opts.ProcessDirectory = debugFunc
	} else {
		opts.ProcessFile = deleteFileFunc
		opts.ProcessDirectory = deleteDirectoryFunc
	}

	opts.Files = flag.Args()
	return
}

func deleteFileFunc(file *model.File) (err error) {
	return os.Remove(file.FullPath)
}

func deleteDirectoryFunc(file *model.File) (err error) {
	return os.RemoveAll(file.FullPath)
}

func debugFunc(file *model.File) (err error) {
	//fmt.Println(file.FullPath)
	return nil
}
