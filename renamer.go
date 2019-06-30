package main

import (
	"bufio"
	"fmt"
	"os"
)

type Rename struct {
	ConfirmRename bool
}

func (r *Rename) Rename(old, new string) (bool, error) {
	if old == new {
		return true, nil
	}

	if r.ConfirmRename {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Rename\n", old, "\nto\n", new, "\n[Y/n]")
		input, _ := reader.ReadString('\n')

		if input != "Y" && input != "y" && input != "\n" {
			// renaming denied by user
			return false, nil
		}
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
