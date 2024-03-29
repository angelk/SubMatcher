package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileInfo struct {
	os.FileInfo
	dir string
}

type fileCollection struct {
	Movies []FileInfo
	Subs   []FileInfo
	Err    error
}

func extractFiles(files []FileInfo) fileCollection {
	movies := make([]FileInfo, 0)
	subs := make([]FileInfo, 0)

	movieExtensions := map[string]bool{
		".avi": true,
		".mkv": true,
		".mp4": true,
		".ts":  true,
	}

	subExtensions := map[string]bool{
		".srt": true,
		".sub": true,
		".sbv": true,
	}

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
func rDirectoryScanner(dir string) (chan []FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileChan := make(chan []FileInfo)
	directoriesToScan := make([]string, 0)

	go func() {
		fileCollection := make([]FileInfo, 0)
		for _, file := range files {
			if file.IsDir() {
				directoriesToScan = append(
					directoriesToScan,
					dir+string(os.PathSeparator)+file.Name(),
				)
			} else {

				fileInfo := FileInfo{
					file,
					dir,
				}
				fileCollection = append(fileCollection, fileInfo)
			}
		}

		fileChan <- fileCollection

		for _, rDir := range directoriesToScan {
			rChan, _ := rDirectoryScanner(rDir)
			for rFileChanData := range rChan {
				fileChan <- rFileChanData
			}
		}

		close(fileChan)
	}()

	return fileChan, nil
}

func directoryScanner(dir string) (chan []FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fileChan := make(chan []FileInfo)

	go func() {
		fileCollection := make([]FileInfo, 0)
		for _, file := range files {
			fileInfo := FileInfo{
				file,
				dir,
			}
			fileCollection = append(fileCollection, fileInfo)
		}

		fileChan <- fileCollection
		close(fileChan)
	}()

	return fileChan, nil
}
