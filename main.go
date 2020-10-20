package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"os"
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
		// make communication channel
		coms := make(chan [2]string)
		defer close(coms)

		// run main function on separate go routine
		go run(*source, *destination, coms)

		// loop to check when data is received over channel
		for {
			// break loop when "done" received over channel
			if e := <- coms; e[0] == "done" {
				break
			} else {
				fmt.Println(e)
			}
		}
	}
}
