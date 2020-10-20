package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"sync"
)

func main() {
	// create argument parser
	parser := argparse.NewParser("File backup", "Copies files from source to destination, ignoring files that haven't been modified")
	// add arguments
	source := parser.String("s", "source", &argparse.Options{Required: true, Help: "The directory to copy files from, only copying files that have been modified"})
	destination := parser.String("d", "destination", &argparse.Options{Required: true, Help: "The directory to copy files to"})

	// run argparse
	if err := parser.Parse(os.Args); err != nil {
		fmt.Println(parser.Usage(err))
	} else {
		// create worker wait group to wait for completion
		var wg sync.WaitGroup

		// run main function separate thread
		wg.Add(1)
		go run(*source, *destination, &wg)

		// run wait group
		wg.Wait()
	}
}
