package model

import (
	"github.com/spf13/afero"
)

// Options allows command line arguments and processing configuration
// to be passed around the application
type Options struct {
	Sweeps           int
	Debug            bool
	Files            []string
	ProcessFile      func(file *File) (err error)
	ProcessDirectory func(file *File) (err error)
}

// NewOptions spins up a pointer to an Options struct
// with some overridable parameters
func NewOptions(fs afero.Fs, sweeps int, debug bool) (opts *Options) {
	opts = new(Options)
	opts.Sweeps = sweeps
	opts.Debug = debug

	if opts.Debug {
		debugFunc := func(file *File) (err error) {
			return nil
		}

		opts.ProcessFile = debugFunc
		opts.ProcessDirectory = debugFunc
	} else {
		opts.ProcessFile = func(file *File) (err error) {
			return fs.Remove(file.FullPath)
		}
		opts.ProcessDirectory = func(file *File) (err error) {
			return fs.RemoveAll(file.FullPath)
		}
	}

	return
}
