package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/akamensky/argparse"
	"github.com/schollz/progressbar/v3"
)

func main() {
	// create argument parser
	parser := argparse.NewParser("GoBack File Backup", "Copies files from source to destination, ignoring files that haven't been modified")
	// add arguments
	source := parser.String("s", "source", &argparse.Options{Required: true, Help: "The directory to copy files from, only copying files that have been modified"})
	destination := parser.String("d", "destination", &argparse.Options{Required: true, Help: "The directory to copy files to"})
	outputCSV := parser.Flag("o", "outputcsv", &argparse.Options{Required: false, Help: "Output data to a CSV file"})

	// run argparse
	if err := parser.Parse(os.Args); err != nil {
		log.Fatal(parser.Usage(err))
	}

	// create main data array
	data := [][]string{}

	// make communication channel
	coms := make(chan []string)

	// run main function on separate go routine
	go run(*source, *destination, coms)

	// create progress bar
	length, _ := strconv.Atoi((<-coms)[1])
	bar := progressbar.Default(int64(length))

	// loop to check when data is received over channel
	for e := range coms {
		data = append(data, e)
		_ = bar.Add(1)
	}

	// create csv file if argument given
	if *outputCSV {
		// create csv file
		file, err := os.Create("output.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// create csv file writer
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
