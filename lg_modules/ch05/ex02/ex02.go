package main

import (
	"fmt"
	"log"
	"os"
)

func fileLen(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return int(info.Size()), nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Need filename.")
	}
	byteNum, err := fileLen(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("File size: %d bytes in file %s\n", byteNum, os.Args[1])
}
