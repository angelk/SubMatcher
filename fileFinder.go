package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
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

var movieExtensions = map[string]bool{
	".avi": true,
	".mkv": true,
	".ts":  true,
}

var subExtensions = map[string]bool{
	".srt": true,
	".sub": true,
	".sbv": true,
}

var incompleteExtension = map[string]bool{
	".!qB": true,
}

type fakeFileInfo struct {
	name string
}

func (f fakeFileInfo) Name() string {
	return f.name
}

func (f fakeFileInfo) Size() int64 {
	return 0
}

func (f fakeFileInfo) Mode() fs.FileMode {
	return 0
}

func (f fakeFileInfo) ModTime() time.Time {
	return time.Now()
}

func (f fakeFileInfo) IsDir() bool {
	return false
}

func (f fakeFileInfo) Sys() any {
	return nil
}

var _ fs.FileInfo = (*fakeFileInfo)(nil)

func extractFiles(files []FileInfo) fileCollection {
	movies := make([]FileInfo, 0)
	subs := make([]FileInfo, 0)

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if _, okIncomp := incompleteExtension[ext]; okIncomp {
			fmt.Println("======== unknownw:", file.Name())
			fmt.Printf("%T [type]---- \n", file)
			/*
				newFile := fakeFileInfo{
					name: file.Name(),
				}
			*/
			// file = newFile
		}

		if _, ok := movieExtensions[ext]; ok {
			movies = append(movies, file)
		} else if _, okSub := subExtensions[ext]; okSub {
			subs = append(subs, file)
		} else {
			fmt.Println("Skipping file", file.Name(), "Unknown extension!")
			fmt.Println("Debug", file.Name, " --- ", ext)
		}
	}

	fc := fileCollection{
		movies,
		subs,
		nil,
	}

	return fc
}

func extractSubsAndVideos(dir string) fileCollection {
	files, error := ioutil.ReadDir(dir)
	if error != nil {
		return fileCollection{
			nil,
			nil,
			error,
		}
	}

	movies := make([]FileInfo, 0)
	subs := make([]FileInfo, 0)

	for _, file := range files {
		ext := filepath.Ext(file.Name())

		if _, ok := movieExtensions[ext]; ok {
			movies = append(movies, FileInfo{
				file,
				dir,
			})
		} else if _, okSub := subExtensions[ext]; okSub {
			subs = append(subs, FileInfo{
				file,
				dir,
			})
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
