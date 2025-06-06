package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

//go:embed rights
var rightsInfo embed.FS

func main() {
	if len(os.Args) == 1 {
		printRightsFiles()
		os.Exit(0)
	}
	data, err := rightsInfo.ReadFile("rights/" + os.Args[1] + "_rights.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func printRightsFiles() {
	fmt.Println("Type a language to choose this translation version of UDHR")
	fmt.Println("languages to choose from:")
	fs.WalkDir(rightsInfo, "rights",
		func(path string, d fs.DirEntry, err error) error {
			if !d.IsDir() {
				_, fileName, _ := strings.Cut(path, "/")
				language, _, _ := strings.Cut(fileName, "_")
				fmt.Println(language)
			}
			return nil
		})
}
