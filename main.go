package kognit

import (
	"fmt"
)

type DirectoryCompressionAlgorithm int
type FileCompressionAlgorithm int
type ImageCompressionAlgorithm int

type CompressionAlgorithm interface {
	Encode(src string) error
	Decode(src string) error
}

const (
	Flate FileCompressionAlgorithm = iota
	Deflate
	Gzip
	Huffman
	LZW
	RLE
)

func (a FileCompressionAlgorithm) Encode(dataPath string) error {
	switch a {
	case Flate:
		fmt.Println("File encoding using Flate")
	case Deflate:
		fmt.Println("File encoding using Deflate")
	case Gzip:
		fmt.Println("File encoding using Gzip")
	case Huffman:
		fmt.Println("File encoding using Huffman")
	case LZW:
		fmt.Println("File encoding using LZW")
	case RLE:
		fmt.Println("File encoding using RLE")
	}
	return nil
}

func (a FileCompressionAlgorithm) Decode(dataPath string) error {
	switch a {
	case Flate:
		fmt.Println("File decoding using Flate")
	case Deflate:
		fmt.Println("File decoding using Deflate")
	case Gzip:
		fmt.Println("File decoding using Gzip")
	case Huffman:
		fmt.Println("File decoding using Huffman")
	case LZW:
		fmt.Println("File decoding using LZW")
	case RLE:
		fmt.Println("File decoding using RLE")
	}
	return nil
}

const (
	JPEG ImageCompressionAlgorithm = iota
	JPEG2000
	PNG
	GIF
)

func (a ImageCompressionAlgorithm) Encode(dataPath string) error {
	switch a {
	case JPEG:
		fmt.Println("File encoding using JPEG")
	case JPEG2000:
		fmt.Println("File encoding using JPEG2000")
	case PNG:
		fmt.Println("File encoding using PNG")
	case GIF:
		fmt.Println("File encoding using GIF")
	}
	return nil
}

func (a ImageCompressionAlgorithm) Decode(dataPath string) error {
	switch a {
	case JPEG:
		fmt.Println("File decoding using JPEG")
	case JPEG2000:
		fmt.Println("File decoding using JPEG2000")
	case PNG:
		fmt.Println("File decoding using PNG")
	case GIF:
		fmt.Println("File decoding using GIF")
	}
	return nil
}
