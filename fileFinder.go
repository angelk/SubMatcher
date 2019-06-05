package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type fileCollection struct {
	Movies []os.FileInfo
	Subs   []os.FileInfo
	Err    error
}

func extractFiles(files []os.FileInfo) fileCollection {
	movies := make([]os.FileInfo, 0)
	subs := make([]os.FileInfo, 0)

	movieExtensions := make(map[string]bool)
	movieExtensions[".avi"] = true
	movieExtensions[".mkv"] = true
	movieExtensions[".ts"] = true

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

	fc := fileCollection{
		movies,
		subs,
		nil,
	}

	return fc
}

func extractSubsAndVideos(dir string, scanners chan []os.FileInfo) fileCollection {
	files, error := ioutil.ReadDir(dir)
	if error != nil {
		return fileCollection{
			nil,
			nil,
			error,
		}
	}

	movies := make([]os.FileInfo, 0)
	subs := make([]os.FileInfo, 0)

	movieExtensions := make(map[string]bool)
	movieExtensions[".avi"] = true
	movieExtensions[".mkv"] = true
	movieExtensions[".ts"] = true

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

	fc := fileCollection{
		movies,
		subs,
		nil,
	}

	return fc
}

// directory scanner used in "-r" flag
func rDirectoryScanner(dir string) (chan []os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileChan := make(chan []os.FileInfo)

	go func() {
		fileCollection := make([]os.FileInfo, 0)
		for _, file := range files {
			fileCollection = append(fileCollection, file)
		}

		fileChan <- fileCollection
		// @TODO add directories for later scanning

		close(fileChan)
	}()

	return fileChan, nil
}

func directoryScanner(dir string) (chan []os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileChan := make(chan []os.FileInfo)

	go func() {
		fileCollection := make([]os.FileInfo, 1)
		for _, file := range files {
			fileCollection = append(fileCollection, file)
		}

		fileChan <- fileCollection
		close(fileChan)
	}()

	return fileChan, nil
}
