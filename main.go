package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	recursiveOption := flag.Bool("r", false, "recursive option")
	noInteraction := flag.Bool("y", false, "No interaction")
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

	renamer := Rename{
		!*noInteraction,
	}

	fmt.Println(renamer)

	for filesCollection := range filesChan {
		matchSubtitles(filesCollection, renamer)
	}
}

func matchSubtitles(files []FileInfo, renamer Rename) {
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

	// remove exact matches
	scanAgain := true
	for scanAgain {
		scanAgain = false
		for movieIndex, movie := range movies {
			movieName := movie.Name()[:len(movie.Name())-len(filepath.Ext(movie.Name()))]
			for subIndex, sub := range subs {
				subName := sub.Name()[:len(sub.Name())-len(filepath.Ext(sub.Name()))]
				if subName == movieName {
					subs[subIndex], subs[len(subs)-1] = subs[len(subs)-1], subs[subIndex]
					subs = subs[:len(subs)-1]

					movies[movieIndex], movies[len(movies)-1] = movies[len(movies)-1], movies[movieIndex]
					movies = movies[:len(movies)-1]
					scanAgain = true
					break
				}
			}
			if scanAgain {
				break
			}
		}
	}

	// Matching
	for _, movie := range movies {
		bestMatchScore := 0
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

		if bestMatchScore == 0 {
			fmt.Println("Skipping score '0'")
			continue
		}

		fmt.Println("score ", bestMatchScore, movie.Name(), bestMatchFile.Name())
		fmt.Println("----")

		movieLenWithoutExt := len(movie.Name()) - len(filepath.Ext(movie.Name()))
		subsExtension := filepath.Ext(bestMatchFile.Name())

		fmt.Println("Matched subs for " + movie.Name())

		renamed, renameError := renamer.Rename(
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
