package lncommon

import (
	"os"
	"strings"
)

/*
./path/to -> /[working directory]/path/to
path/to   -> path/to
/path/to  -> /path/to
*/
func ToPath(path *string) string {
	// set resource path
	var truePath string
	if strings.HasPrefix(truePath, "./") {
		// get current directory
		currentDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		truePath = currentDir + (*path)[1:]
	} else {
		truePath = *path
	}

	return truePath
}

func Copy(srcPath *string, destPath *string) {
	bytesRead, err := os.ReadFile(*srcPath)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(*destPath, bytesRead, 0644)
	if err != nil {
		panic(err)
	}
}
