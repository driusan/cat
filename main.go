package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var unbuffered, halp bool

func main() {
	flag.BoolVar(&unbuffered, "u", false, "Do not buffer output")
	flag.BoolVar(&halp, "h", false, "Print help")
	flag.Parse()
	filenames := flag.Args()

	if halp {
		fmt.Printf("Usage: %s [flags] [file]\n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}

	var erred bool
	if len(filenames) == 0 {
		if err := CatFile("-"); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			erred = true
		}
	} else {
		for _, filename := range filenames {
			if err := CatFile(filename); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				erred = true
			}
		}
	}
	if erred {
		os.Exit(2)
	}
}

func CatFile(filename string) error {
	var f *os.File
	var err error

	if filename == "-" {
		f = os.Stdin
	} else {
		f, err = os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	tee := io.TeeReader(f, os.Stdout)

	var buf []byte
	if unbuffered {
		buf = make([]byte, 1)
	} else {
		buf = make([]byte, 1024)
	}
	for _, err := tee.Read(buf); err != io.EOF; _, err = tee.Read(buf) {
		_, err = tee.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
	}
	return nil
}
