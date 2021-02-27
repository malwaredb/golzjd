package main

import (
	"github.com/malwaredb/golzjd"
	"os"
	"fmt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usages:\n")
		fmt.Fprintf(os.Stderr, "\t%s <File> to see LZJD hash of a file\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\t%s <File> <File> to use LZJD to calculate similarity\n", os.Args[0])
		os.Exit(1)
	}

	firstFile := golzjd.GenerateHashFromFile(os.Args[1])
	if len(os.Args) > 2 {
		secondFile := golzjd.GenerateHashFromFile(os.Args[2])
		similarity := golzjd.CompareHashes(firstFile, secondFile)
		fmt.Printf("Similarity: %d\n", similarity)
	} else {
		fmt.Printf("Hash: %s\n", firstFile)
	}
}
