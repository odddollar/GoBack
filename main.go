package main

import (
	"encoding/csv"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/olekukonko/tablewriter"
	"github.com/schollz/progressbar/v3"
	"log"
	"os"
	"strconv"
)

func main() {
	// create argument parser
	parser := argparse.NewParser("File backup", "Copies files from source to destination, ignoring files that haven't been modified")
	// add arguments
	source := parser.String("s", "source", &argparse.Options{Required: true, Help: "The directory to copy files from, only copying files that have been modified"})
	destination := parser.String("d", "destination", &argparse.Options{Required: true, Help: "The directory to copy files to"})
	printTable := parser.Flag("p", "printtable", &argparse.Options{Required: false, Help: "Print an output table with the action taken for each file found in source directory"})
	outputCSV := parser.Flag("o", "outputcsv", &argparse.Options{Required: false, Help: "Output data to CSV file"})

	// run argparse
	if err := parser.Parse(os.Args); err != nil {
		fmt.Println(parser.Usage(err))
	} else {
		// create main data array
		data := [][]string{}

		// make communication channel
		coms := make(chan []string)

		// run main function on separate go routine
		go run(*source, *destination, coms)

		// create table writer
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"File", "Status"})
		table.SetRowLine(true)
		table.SetAutoWrapText(false)

		// create progress bar
		temp := <- coms
		length, _ := strconv.Atoi(temp[1])
		bar := progressbar.Default(int64(length))

		// loop to check when data is received over channel
		for e := range coms {
			table.Append(e)
			data = append(data, e)
			_ = bar.Add(1)
		}

		// print table if argument given
		if *printTable {
			table.Render()
		}

		// create csv file if argument given
		if *outputCSV {
			// create csv file
			file, err := os.Create("output.csv")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			// create csv file
			writer := csv.NewWriter(file)
			defer writer.Flush()

			// write header
			_ = writer.Write([]string{"File", "Status"})

			// write remaining data
			for _, row := range data {
				_ = writer.Write(row)
			}
		}
	}
}
