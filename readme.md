# golzjd

This project is a work in progress to convert Edward Raff's original [LZJD](https://github.com/EdwardRaff/LZJD) to Go. To use LZJD with your project, run `go get github.com/malwaredb/golzjd`, and import the project into your Go file. For more information about LZJD itself, read the [original paper](https://arxiv.org/abs/1708.03346).

## Example use

* Get the LZJD hash of a file: `lzjdHash := golzjd.GenerateHashFromFile("/path/to/file")`
* Compute the similarity of two files: `similarity := golzjd.CompareHashes(hash1, hash2)`

## Requirements

* BOOST libraries
* C++11 compiler (likely gcc)

This code has been tested on Ubuntu 16.04 with Go 1.10 & gcc 5.4.

## Notice

This project includes files from the original [LZJD](https://github.com/EdwardRaff/LZJD). Refer to that project page for any updates.