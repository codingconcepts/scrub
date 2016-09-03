package model

// Options allows command line arguments and processing configuration
// to be passed around the application
type Options struct {
	Sweeps           int
	Debug            bool
	Files            []string
	ProcessFile      func(file *File) (err error)
	ProcessDirectory func(file *File) (err error)
}
