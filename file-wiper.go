/*
	TODO: move the FileSystem somewhere else so i dont have to inject it everywhere
	TODO: make options creation a bit less shit
*/
package main

import (
	"flag"
	"fmt"

	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/bassrob/file-wiper/helper"
	"github.com/bassrob/file-wiper/model"
	"github.com/spf13/afero"
)

var (
	FileSystem  afero.Fs
	ProgressBar *pb.ProgressBar
	Options     *model.Options
)

func main() {
	Options = parse()
	FileSystem = afero.NewOsFs()

	if err := processFiles(); err != nil {
		fmt.Printf("An error occurred processing files: %s", err)
	}

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
	ProgressBar = createProgressBar(totalSize)

	// loop through all of the nested files and process them
	ProgressBar.Start()
	for _, file := range files {
		if err = helper.ProcessFile(FileSystem, Options, file); err != nil {
			panic(err)
		}
		ProgressBar.Add64(file.Size)
	}
	ProgressBar.Finish()

	return
}

func processDirectories() (err error) {
	var directories []*model.File
	if directories, err = helper.GetTopLevelDirectories(FileSystem, Options.Files); err != nil {
		return
	}

	// now that all of the nested files have been cleaned up,
	// remove the top-level directories
	for _, file := range directories {
		if err := Options.ProcessDirectory(file); err != nil {
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
