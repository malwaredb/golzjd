package golzjd

// #cgo CXXFLAGS: -Wall -std=c++11 -I .
// #cgo LDFLAGS: -lstdc++
/*
#include "lzjd_helper.h"
*/
import "C"
import (
	"unsafe"
)

func CompareHashes(hash1 string, hash2 string) int32 {
	CHash1 := C.CString(hash1)
	CHash2 := C.CString(hash2)
	result := C.lzjd_similarity(CHash1, CHash2)
	return int32(result)
}

func GenerateHashFromFile(file string) string {
	Cfile := C.CString(file)
	result := C.createDigest(Cfile)
	return C.GoString(result)
}

func GenerateHashFromBuffer(buffer []byte) string {
	cBuffArray := (*C.char)(unsafe.Pointer(&buffer[0]))
	result := C.createDigestFromBuffer(cBuffArray, C.int(len(buffer)))
	return C.GoString(result)
}
