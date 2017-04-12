package main

import (
	"flag"
	"fmt"
	"log"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/codingconcepts/file-wiper/helper"
	"github.com/codingconcepts/file-wiper/model"
	"github.com/spf13/afero"
)

var (
	// FileSystem allows for the mocking and easier
	// testing of the underlying file system
	FileSystem afero.Fs

	// Options allows the basic processing behaviour
	// to be overriden
	Options *model.Options
)

func main() {
	Options = parse()
	FileSystem = afero.NewOsFs()

	// overwrite every file with pseudo random data
	if err := processFiles(); err != nil {
		fmt.Printf("An error occurred processing files: %s", err)
	}

	// delete the top-level directories containing the files (if any)
	if err := processDirectories(); err != nil {
		fmt.Printf("An error occurred processing directories: %s", err)
	}

	printAndHold("Finished, press any key to quit")
}

func processFiles() (err error) {
	var files []*model.File
	if files, err = helper.GetAllFiles(FileSystem, Options.Files); err != nil {
		return
	}

	totalSize := model.TotalSize(files)
	progressBar := createProgressBar(totalSize)

	progressBar.Start()
	for _, file := range files {
		if err = helper.ProcessFile(FileSystem, Options, file); err != nil {
			log.Fatal(err)
		}
		progressBar.Add64(file.Size)
	}
	progressBar.Finish()

	return
}

func processDirectories() (err error) {
	var directories []*model.File
	if directories, err = helper.GetTopLevelDirectories(FileSystem, Options.Files); err != nil {
		return
	}

	for _, file := range directories {
		if err := Options.ProcessDirectory(file); err != nil {
			log.Fatal(err)
		}
	}

	return
}

func createProgressBar(totalSize int64) (bar *pb.ProgressBar) {
	bar = pb.New64(totalSize)
	bar.SetUnits(pb.U_BYTES)

	return
}

func printAndHold(msg string) {
	fmt.Println(msg)
	fmt.Scanln()
}

func parse() (opts *model.Options) {
	var sweeps int
	var debug bool
	flag.IntVar(&sweeps, "s", 10, "the number of overwrite sweeps")
	flag.BoolVar(&debug, "d", false, "set flag to print files/directories, not delete them")
	flag.Parse()

	opts = model.NewOptions(FileSystem, sweeps, debug)
	opts.Files = flag.Args()
	return
}
