package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getMatchScore(a string, b string, diffLen int) (int, error) {
	// our algorithm doesn't work with short names
	if len(a) < 6 || len(b) < 6 {
		return 0, errors.New("Require string with len >= 6")
	}

	matchScore := 0

	aLen := len(a)
	bLen := len(b)
	for i := 0; i+diffLen-1 < aLen; i++ {
		aPart := a[i : i+diffLen]
		for j := 0; j+diffLen-1 < bLen; j++ {
			bPart := b[j : j+diffLen]
			if aPart == bPart {
				matchScore++
			}
		}
	}

	return matchScore, nil
}

func extractFiles(dir string) ([]os.FileInfo, []os.FileInfo, error) {
	files, error := ioutil.ReadDir(dir)
	if error != nil {
		return nil, nil, error
	}

	movies := make([]os.FileInfo, 0)
	subs := make([]os.FileInfo, 0)

	movieExtensions := make(map[string]bool)
	movieExtensions[".avi"] = true
	movieExtensions[".mkv"] = true

	subExtensions := make(map[string]bool)
	subExtensions[".srt"] = true
	subExtensions[".sub"] = true
	subExtensions[".sbv"] = true

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if _, ok := movieExtensions[ext]; ok {
			movies = append(movies, file)
		} else if _, okSub := subExtensions[ext]; okSub {
			subs = append(subs, file)
		} else {
			fmt.Println("Skipping file", file.Name(), "Unknown extension!")
		}
	}

	return movies, subs, nil
}

func main() {
	directory := "/home/potaka/Projects/movie/testCases/"
	movies, subs, _ := extractFiles(directory)

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
		var bestMatchFile os.FileInfo
		var bestMatchIndex int

		for subIndex, sub := range subs {
			tempMatchScore, _ := getMatchScore(movie.Name(), sub.Name(), 3)

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

		// probably we should ask for confirmation!
		movieLenWithoutExt := len(movie.Name()) - len(filepath.Ext(movie.Name()))
		subsExtension := filepath.Ext(bestMatchFile.Name())

		// rename subs file
		os.Rename(
			directory+bestMatchFile.Name(),
			directory+movie.Name()[0:movieLenWithoutExt]+subsExtension,
		)

		// remove subs from the list
		subs[bestMatchIndex], subs[len(subs)-1] = subs[len(subs)-1], subs[bestMatchIndex]
		subs = subs[:len(subs)-1]
	}
}
