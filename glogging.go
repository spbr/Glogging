package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

//Config is a structure that will hold all the command line arguments
type Config struct {
	SrcDir      string
	DestDir     string
	Concurrency int
}

var config = Config{}

func init() {
	flag.StringVar(&config.SrcDir, "src", "", "source directory")
	flag.StringVar(&config.DestDir, "dest", "", "destination directory")
	flag.IntVar(&config.Concurrency, "concurrency", 1, "concurrency")
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

	if config.SrcDir == config.DestDir {
		fmt.Println("Source and Destination cannot match")
		os.Exit(64)
	}
	finfo, err := os.Stat(config.SrcDir)
	if err != nil {
		fmt.Printf("Unable to validate source directory: %s\n", err.Error())
		os.Exit(64)
	}
	if !finfo.IsDir() {
		fmt.Println("Source is not a directory")
		os.Exit(2)
	}
	finfo, err = os.Stat(config.DestDir)
	if err != nil {
		fmt.Printf("Unable to validate destination directory: %s\n", err.Error())
		os.Exit(64)
	}
	if !finfo.IsDir() {
		fmt.Println("Destination is not a directory")
		os.Exit(2)
	}

	if !strings.HasSuffix(config.DestDir, string(os.PathSeparator)) {
		config.DestDir = config.DestDir + string(os.PathSeparator)
	}
	if !strings.HasSuffix(config.SrcDir, string(os.PathSeparator)) {
		config.SrcDir = config.SrcDir + string(os.PathSeparator)
	}
	/*
	 ** Let's spin through directory
	 */

	filechannel := make(chan os.FileInfo)
	var wg sync.WaitGroup
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go archiveFiles(&wg, filechannel, &config)
	}

	fileinfos, err := ioutil.ReadDir(config.SrcDir)
	if err != nil {
		fmt.Printf("Unable to read source directory: %s\n", err.Error())
		os.Exit(64)
	}
	for _, file := range fileinfos {
		if file.IsDir() {
			continue
		}
		if isHidden(file.Name()) {
			continue
		} else {
			filechannel <- file
		}

	}

	close(filechannel)
	wg.Wait()
	os.Exit(0)
}

// isHidden determins if the file is hidden
func isHidden(filename string) bool {
	if strings.HasPrefix(filename, ".") {
		return true
	}
	return false
}

//archiveFiles takes a file and moves it to the destination
func archiveFiles(wg *sync.WaitGroup, filechannel chan os.FileInfo, config *Config) {
	defer wg.Done()
	for {
		file, more := <-filechannel
		if !more {
			return
		}
		_, _ = copyFile(config, file)
	}
}

func copyFile(config *Config, file os.FileInfo) (bool, error) {
	fullpath := fmt.Sprintf("%s%s", config.SrcDir, file.Name())
	tempDestPath := fmt.Sprintf("%s.%s", config.DestDir, file.Name())
	destPath := fmt.Sprintf("%s%s", config.DestDir, file.Name())
	ifile, err := os.Open(fullpath)
	if err != nil {
		return false, err
	}
	defer ifile.Close()

	dfile, err := os.Create(tempDestPath)
	if err != nil {
		return false, err
	}

	data := make([]byte, 4096000)
	for {
		data = data[:cap(data)]
		n, err := ifile.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading in file: %s\n", err.Error())
			dfile.Close()
			os.Remove(tempDestPath)
			return false, err
		}
		data = data[:n]
		dfile.Write(data)
	}
	dfile.Close()
	err = os.Rename(tempDestPath, destPath)
	if err != nil {
		os.Remove(tempDestPath)
		return false, err
	}
	err = os.Chmod(destPath, file.Mode())
	if err != nil {
		os.Remove(destPath)
		return false, err
	}
	return true, nil
}
