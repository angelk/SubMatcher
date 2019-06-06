package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func rename(old, new string) (bool, error) {
	if old == new {
		return true, nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Rename\n", old, "\nto\n", new, "\n[Y/n]")
	input, _ := reader.ReadString('\n')

	if input != "Y" && input != "y" && input != "\n" {
		// renaming denied by user
		return false, nil
	}

	// rename subs file
	error := os.Rename(
		old,
		new,
	)

	if error != nil {
		return false, error
	}

	return true, nil
}

func main() {
	recursiveOption := flag.Bool("r", false, "recursive option")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		log.Fatalln("Error, path (1st argument) not provided")
	}

	directory := args[0]
	directory = strings.TrimRight(directory, string(os.PathSeparator))

	var filesChan chan []FileInfo
	var err error

	if *recursiveOption {
		filesChan, err = rDirectoryScanner(directory)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		filesChan, err = directoryScanner(directory)
		if err != nil {
			log.Fatalln(err)
		}
	}

	for filesCollection := range filesChan {
		matchSubtibles(filesCollection)
	}

}

func matchSubtibles(files []FileInfo) {
	fc := extractFiles(files)

	movies := fc.Movies
	subs := fc.Subs
	extractFilesError := fc.Err

	if extractFilesError != nil {
		log.Fatalln(extractFilesError)
	}

	fmt.Println("--- Movies")
	for _, file := range movies {
		fmt.Println(file.Name())
	}

	fmt.Println("--- Subs")
	for index, file := range subs {
		fmt.Println(index, file.Name())
	}

	fmt.Println("--- --- ---")

	// Matching
	for _, movie := range movies {
		var bestMatchScore int
		var bestMatchFile FileInfo
		var bestMatchIndex int

		for subIndex, sub := range subs {
			tempMatchScore := getMatchScore(movie.Name(), sub.Name(), 3)

			// @TODO check for same score!
			if bestMatchScore < tempMatchScore {
				bestMatchScore = tempMatchScore
				bestMatchFile = sub
				bestMatchIndex = subIndex
			}
		}

		fmt.Println("score ", bestMatchScore, movie.Name(), bestMatchFile.Name())
		fmt.Println("----")

		if bestMatchScore == 0 {
			fmt.Println("Skipping score '0'")
			continue
		}

		movieLenWithoutExt := len(movie.Name()) - len(filepath.Ext(movie.Name()))
		subsExtension := filepath.Ext(bestMatchFile.Name())

		fmt.Println("Matched subs for " + movie.Name())

		renamed, renameError := rename(
			bestMatchFile.dir+string(os.PathSeparator)+bestMatchFile.Name(),
			movie.dir+string(os.PathSeparator)+movie.Name()[0:movieLenWithoutExt]+subsExtension,
		)

		if renameError != nil {
			fmt.Println(renameError)
		}

		if renamed {
			// remove subs from the list
			subs[bestMatchIndex], subs[len(subs)-1] = subs[len(subs)-1], subs[bestMatchIndex]
			subs = subs[:len(subs)-1]
		}
	}
}
