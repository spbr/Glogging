package main

import (
	"flag"
	"fmt"
	"os"
)

var srcDir string
var destDir string
var concurrency int

func init() {
	flag.StringVar(&srcDir, "src", "", "source directory")
	flag.StringVar(&destDir, "dest", "", "destination directory")
	flag.IntVar(&concurrency, "concurrency", 1, "concurrency")
}

func main() {

	/*
	 ** Parse arguments
	 */
	flag.Usage = func() {
		flag.PrintDefaults()
		os.Exit(64)
	}
	flag.Parse()

	if flag.Lookup("src").Value.String() == "" {
		f := flag.Lookup("src")
		fmt.Printf("Missing argument %s : %s\n", f.Name, f.Usage)
		os.Exit(1)
	}

	if flag.Lookup("dest").Value.String() == "" {
		f := flag.Lookup("dest")
		fmt.Printf("Missing argument %s : %s\n", f.Name, f.Usage)
		os.Exit(1)
	}

	/*
	 ** Validate src and destination
	 */

	if srcDir == destDir {
		fmt.Println("Source and Destination cannot match")
		os.Exit(64)
	}
	finfo, err := os.Stat(srcDir)
	if err != nil {
		fmt.Printf("Unable to validate source directory: %s\n", err.Error())
		os.Exit(64)
	}
	if !finfo.IsDir() {
		fmt.Println("Source is not a directory")
		os.Exit(2)
	}
	finfo, err = os.Stat(destDir)
	if err != nil {
		fmt.Printf("Unable to validate destination directory: %s\n", err.Error())
		os.Exit(64)
	}
	if !finfo.IsDir() {
		fmt.Println("Destination is not a directory")
		os.Exit(2)
	}
}
