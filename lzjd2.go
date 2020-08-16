package golzjd

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

const digest_size uint64 = 1024

type Int32Array []int32

func (f Int32Array) Len() int {
	return len(f)
}

func (f Int32Array) Less(i, j int) bool {
	return f[i] < f[j]
}

func (f Int32Array) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func GetAllHashes(data []byte) []int32 {
	ints := make([]int32, 0)

	x_set := make(map[int32]bool)
	running_hash := NewMurmurHash3()

	for _, b := range data {
		hash := running_hash.PushByte(int8(b))
		//if idx < 2 {
		//	fmt.Printf("Go MurmurHash(%x) == %d\n", int8(b), hash)
		//}
		_, okay := x_set[hash]
		if !okay {
			x_set[hash] = true
			ints = append(ints, hash)
			running_hash.Reset()
		}
	}

	return ints
}

func Digest(k uint64, data []byte) []int32 {
	ints := Int32Array(GetAllHashes(data))
	sort.Sort(ints)
	if uint64(len(ints)) > k {
		ints = ints[:k]
	}
	//fmt.Printf(" Go Ints:\n")
	//fmt.Print(ints)
	//fmt.Print("\n")

	return ints
}

func IntersectVector(left, right []int32) int {
	count := 0

	for _, x := range left {
		for _, y := range right {
			if x == y {
				count += 1
			}
		}
	}

	return count
}

func Similarity (left, right []int32) int {
	same := float64(IntersectVector(left, right))
	sim := same / ( float64(len(left)) + float64(len(right)) - same)
	return int(100.0*sim)
}

func GetLZJDHashFromFile(fname string) string {
	fileContents, err := ioutil.ReadFile(fname)
	if err != nil {
		return ""
	}
	return GetLZJDHashFromBytes(filepath.Base(fname), fileContents)
}

func GetLZJDHashFromBytes(fname string, fcontents []byte) string {
	ints := Digest(digest_size, fcontents)
	thebytes := make([]byte, 0)
	for _, intval := range ints {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(intval))
		for _, bv := range bs {
			thebytes = append(thebytes, bv)
		}
	}
	return "lzjd:"+fname+":"+base64.StdEncoding.EncodeToString(thebytes)
}

func LZJDHashToInts(lzjdHash string) []int32 {
	parts := strings.Split(lzjdHash, ":")
	base64String := parts[len(parts)-1]
	if !strings.HasSuffix(base64String, "==") {
		base64String = base64String + "=="
	}
	lzjdBytes, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		fmt.Printf("String: %s\n", base64String)
		fmt.Printf("Error converting from base64: %s\n", err)
		return nil
	}
	theInts := make([]int32,0)
	for offset := 0; offset < len(lzjdBytes); offset += 4 {
		segmentBytes := []byte{lzjdBytes[offset], lzjdBytes[offset+1], lzjdBytes[offset+2], lzjdBytes[offset+3]}
		thisInt := binary.LittleEndian.Uint32(segmentBytes)
		theInts = append(theInts, int32(thisInt))
	}
	return theInts
}

func CompareHashesPureGo(left, right string) int {
	intsOne := LZJDHashToInts(left)
	if intsOne == nil {
		fmt.Printf("Error converting left hash to ints.\n")
		return -1
	}
	intsTwo := LZJDHashToInts(right)
	if intsTwo == nil {
		fmt.Printf("Error converting right hash to ints.\n")
		return -1
	}
	return Similarity(intsOne, intsTwo)
}