package main

import(
	"io/ioutil"
	"os"
	"fmt"
	"path/filepath"
)

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

	return movies, subs, nil
}
