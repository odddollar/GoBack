package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/olekukonko/tablewriter"
	"os"
)

func main() {
	// create argument parser
	parser := argparse.NewParser("File backup", "Copies files from source to destination, ignoring files that haven't been modified")
	// add arguments
	source := parser.String("s", "source", &argparse.Options{Required: true, Help: "The directory to copy files from, only copying files that have been modified"})
	destination := parser.String("d", "destination", &argparse.Options{Required: true, Help: "The directory to copy files to"})
	printTable := parser.Flag("t", "printtable", &argparse.Options{Required: false, Help: "Print an output table with the action taken for each file found in source directory"})

	// run argparse
	if err := parser.Parse(os.Args); err != nil {
		fmt.Println(parser.Usage(err))
	} else {
		// make communication channel
		coms := make(chan []string)
		defer close(coms)

		// run main function on separate go routine
		go run(*source, *destination, coms)

		// create table writer
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"File", "Status"})
		table.SetRowLine(true)
		table.SetAutoWrapText(false)

		// loop to check when data is received over channel
		for {
			// break loop when "done" received over channel
			if e := <- coms; e[0] == "done" {
				break
			} else {
				table.Append(e)
			}
		}

		if *printTable {
			table.Render()
		}
	}
}
