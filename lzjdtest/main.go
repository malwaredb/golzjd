package main

import (
	"fmt"
	"github.com/malwaredb/golzjd"
	"os"
	"strconv"
	"strings"
	"math/rand"
	"time"
)

func CompareOne(theBin string) {
	raffLZJD := golzjd.GenerateHashFromFile(theBin)
	zakLZJD := golzjd.GetLZJDHashFromFile(theBin)
	fmt.Printf("C++: %s\n", raffLZJD)
	fmt.Printf("Pure Go: %s\n", zakLZJD)
	if strings.Compare(raffLZJD, zakLZJD) == 0 {
		fmt.Println("Same!")
	} else {
		fmt.Println("Not equal.")
	}

	raffSimilarity := golzjd.CompareHashes(raffLZJD, zakLZJD)
	fmt.Printf("C++ similarity: %d\n", raffSimilarity)

	zakSimilarity := golzjd.CompareHashesPureGo(raffLZJD, zakLZJD)
	fmt.Printf("Pure Go similarity: %d\n", zakSimilarity)
}

func CompareTwo(left, right string) {
	raffLZJDLeft := golzjd.GenerateHashFromFile(left)
	zakLZJDLeft := golzjd.GetLZJDHashFromFile(left)

	raffLZJDRight := golzjd.GenerateHashFromFile(right)
	zakLZJDRight := golzjd.GetLZJDHashFromFile(right)

	raffSimilarity := golzjd.CompareHashes(raffLZJDLeft, raffLZJDRight)
	fmt.Printf("C++ similarity: %d\n", raffSimilarity)

	zakSimilarity := golzjd.CompareHashesPureGo(zakLZJDLeft, zakLZJDRight)
	fmt.Printf("Pure Go similarity: %d\n", zakSimilarity)
}

func CompareRandom(rounds int) {
	rand.Seed(time.Now().UnixNano())
	for index := 0; index < rounds; index++ {
		byteArray := make([]byte, 10240)
		rand.Read(byteArray)
		raffLZJD := golzjd.GenerateHashFromBuffer(byteArray)
		zakLZJD := golzjd.GetLZJDHashFromBytes("buffer", byteArray)

		if strings.Compare(raffLZJD, zakLZJD) != 0 {
			fmt.Printf("LZJD mismatch iteration %d\n", index)
			fmt.Printf("%s\n%s\n", raffLZJD, zakLZJD)
			return
		}
	}
}

func main() {
	if len(os.Args) == 3 {
		CompareTwo(os.Args[1], os.Args[2])
		return
	}

	if len(os.Args) == 2 {
		iters, err := strconv.Atoi(os.Args[1])
		if err == nil {
			CompareRandom(iters)
			return
		}
	}

	theBin := os.Args[0]
	if len(os.Args) == 2 {
		theBin = os.Args[1]
		info, err := os.Stat(theBin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if info.IsDir() {
			fmt.Fprintf(os.Stderr, "%s isn't a file.\n", theBin)
			os.Exit(1)
		}
	}
	CompareOne(theBin)
}