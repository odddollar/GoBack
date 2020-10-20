package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func run(source, destination string, coms chan []string) {
	// get list of folders and files in source directory
	sourceFolders, sourceFiles := listFoldersFiles(source)

	// send initial length data over channel
	coms <- []string{"", strconv.Itoa(len(sourceFiles))}

	// remove root from sourceFolders and sourceFiles
	for x := 0; x < len(sourceFolders); x++ {
		sourceFolders[x] = strings.Replace(sourceFolders[x], source, "", 1)
	}
	for x := 0; x < len(sourceFiles); x++ {
		sourceFiles[x] = strings.Replace(sourceFiles[x], source, "", 1)
	}

	// recreate folder structure, creating folders in destination if they don't exist
	recreateFolderStructure(sourceFolders, destination)

	// copy files from source to destination
	for x := 0; x < len(sourceFiles); x++ {
		testFileSource := source + sourceFiles[x]

		// create new filename with appended modification date
		testFileExtension := filepath.Ext(sourceFiles[x])
		testFileName := sourceFiles[x][0:len(sourceFiles[x])-len(testFileExtension)]
		testFileModificationDate := appendModificationDate(testFileSource)
		testFile := destination + testFileName + " " + testFileModificationDate + testFileExtension

		// check if file exists in destination, copy if it doesn't
		if t, _ := exists(testFile); t == true {
			// create temporary array to send over channel
			temp := []string{testFile, "Not copied"}
			// send completed file data over channel
			coms <- temp
		} else {
			// create temporary array to send over channel
			temp := []string{testFile, "Copied"}
			// file doesn't exist, copy across
			_ = copyFiles(testFileSource, testFile)
			// send completed file data over channel
			coms <- temp
		}
	}

	// close com channel
	close(coms)
}

func appendModificationDate(file string) string {
	// get modification date/time and convert to string
	stats, _ := os.Stat(file)
	modTime := stats.ModTime()
	modTimeString := modTime.String()

	// split to remove second part
	modTimeString = strings.Split(modTimeString, ".")[0]

	if strings.Contains(modTimeString, "+") {
		modTimeString = strings.Split(modTimeString, "+")[0]
	}

	// replace bad characters
	modTimeString = strings.ReplaceAll(modTimeString, "-", "")
	modTimeString = strings.ReplaceAll(modTimeString, ":", "")
	modTimeString = strings.ReplaceAll(modTimeString, " ", "")

	return modTimeString
}

func copyFiles(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func recreateFolderStructure(sourceFolders []string, destination string) {
	// iterate through source folders, checking if they exist to recreate file structure
	for x := 0; x < len(sourceFolders); x++ {
		testDirectory := destination + sourceFolders[x]
		// check if directory exists in destination folder
		if t, _ := exists(testDirectory); t == false {
			_ = os.Mkdir(testDirectory, 0777)
		}
	}
}

func listFoldersFiles(root string) ([]string, []string) {
	// create slices to return
	var folders []string
	var files []string

	// walk through directories and folders
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {return err}

		// append data to relevant slice based on dot in file extension
		if !strings.Contains(path, ".") {
			folders = append(folders, path)
		} else {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return folders, files
}

// check if folder exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}