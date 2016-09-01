/*
	// TODO: bang err creation into if blocks
	// TODO: process directories
	// TODO: progress indication
*/

package main

import (
	"bufio"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"os"
)

type options struct {
	sweeps int
	files  []string
}

func main() {
	opts := parse()
	processFiles(opts)

	printAndHold("Done")
}

func processFiles(opts *options) (err error) {
	for _, file := range opts.files {
		// overwrite the file with random data
		for i := 0; i < opts.sweeps; i++ {
			if err = overwrite(file); err != nil {
				return err
			}
		}

		// delete the file
		if err = os.Remove(file); err != nil {
			return err
		}
	}

	return
}

func overwrite(filePath string) (err error) {
	writer, err := createWriter(filePath)
	if err != nil {
		return
	}

	stats, err := os.Stat(filePath)
	if err != nil {
		return
	}

	// using a limited reader so we don't generate unlimited data
	limitReader := io.LimitReader(rand.Reader, stats.Size())
	return pipe(limitReader, writer)
}

func createWriter(filePath string) (writer *bufio.Writer, err error) {
	outputFile, err := os.OpenFile(filePath, os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return bufio.NewWriter(outputFile), nil
}

func pipe(reader io.Reader, writer *bufio.Writer) (err error) {
	_, err = writer.ReadFrom(reader)
	if err != nil {
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
	flag.Parse()

	opts.files = flag.Args()
	return
}
