package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/csothen/kognit/pkg/compress"
)

const (
	lz77    = "LZ77"
	lzw     = "LZW"
	rle     = "RLE"
	huffman = "Huffman"
)

// Compress handles the compression of files making calls to different
// compression algorithms
func Compress(filePath string, algorithm string) error {

	// Read file content and check if the reading didn't fail
	fileContent, rErr := ioutil.ReadFile(filePath)
	if rErr != nil {
		log.Fatal(rErr)
		return rErr
	}

	// Variable to hold the compressed file
	var compressedData []byte
	// Variable to hold the compression error (if any)
	var cErr error

	// Compress the file using the wanted algorithm and
	// check if there was an error during the compression process
	switch algorithm {
	case lz77:
		// Compress the file using the LZ77 compression algorithm
		compressedData, cErr = compress.LZ77(fileContent)
		break
	case lzw:
		// Compress the file using the LZW compression algorithm
		compressedData, cErr = compress.LZW(fileContent)
		break
	case rle:
		// Compress the file using the RLE compression algorithm
		compressedData, cErr = compress.RLE(fileContent)
		break
	default:
		// Compress the file using the Huffman compression algorithm
		compressedData, cErr = compress.Huffman(fileContent)
		break
	}

	if cErr != nil {
		log.Fatal(cErr)
		return cErr
	}

	// Write compressed content and check if writting didnt fail

	// Input file's location
	location := path.Dir(filePath)
	// Output file's name
	name := strings.Split(path.Base(filePath), ".")[0] + ".kgi"

	// Location of the new compressed file
	newFile := path.Join(location, name)

	err := ioutil.WriteFile(newFile, compressedData, 6666)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}

	return nil
}
