package kognit

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ItemType int
type DirectoryCompressionAlgorithm int
type FileCompressionAlgorithm int
type ImageCompressionAlgorithm int

type CompressionAlgorithm interface {
	Encode(dataPath string) error
	Decode(dataPath string) error
}

const (
	Directory ItemType = iota
	File
	Image
)

const (
	ZIP DirectoryCompressionAlgorithm = iota
	TAR
)

func (a DirectoryCompressionAlgorithm) Encode(dataPath string) error {
	output := dataPath

	switch a {
	case ZIP:
		output += ".zip"
		if err := encodeZipArchive(dataPath, output); err != nil {
			return err
		}
	case TAR:
		output += ".tar"
		if err := encodeTarArchive(dataPath, output); err != nil {
			return err
		}
	default:
		return errors.New("Compression algorithm does not exist")
	}
	return nil
}

func encodeZipArchive(dataPath, output string) error {
	files, err := allDirFiles(dataPath)
	if err != nil {
		return err
	}

	zipFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err = addToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addToZip(w *zip.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = filename
	header.Method = zip.Deflate

	writer, err := w.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func encodeTarArchive(dataPath, output string) error {
	files, err := allDirFiles(dataPath)
	if err != nil {
		return nil
	}

	tarFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	gzWriter := gzip.NewWriter(tarFile)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	for _, file := range files {
		if err = addToTar(tarWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func addToTar(w *tar.Writer, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = filename

	if err := w.WriteHeader(header); err != nil {
		return err
	}

	_, err = io.Copy(w, file)
	return err
}

func allDirFiles(dirPath string) ([]string, error) {
	files := []string{}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

func (a DirectoryCompressionAlgorithm) Decode(dataPath string) error {
	switch a {
	case ZIP:
	case TAR:
	default:
		return errors.New("Compression algorithm does not exist")
	}
	return nil
}

const (
	Flate FileCompressionAlgorithm = iota
	Gzip
	Huffman
	LZW
	RLE
)

func (a FileCompressionAlgorithm) Encode(dataPath string) error {
	switch a {
	case Flate:
		fmt.Println("File encoding using Flate")
	case Gzip:
		fmt.Println("File encoding using Gzip")
	case Huffman:
		fmt.Println("File encoding using Huffman")
	case LZW:
		fmt.Println("File encoding using LZW")
	case RLE:
		fmt.Println("File encoding using RLE")
	default:
		return errors.New("Compression algorithm does not exist")
	}
	return nil
}

func (a FileCompressionAlgorithm) Decode(dataPath string) error {
	switch a {
	case Flate:
		fmt.Println("File decoding using Flate")
	case Gzip:
		fmt.Println("File decoding using Gzip")
	case Huffman:
		fmt.Println("File decoding using Huffman")
	case LZW:
		fmt.Println("File decoding using LZW")
	case RLE:
		fmt.Println("File decoding using RLE")
	default:
		return errors.New("Compression algorithm does not exist")
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
	default:
		return errors.New("Compression algorithm does not exist")
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
	default:
		return errors.New("Compression algorithm does not exist")
	}
	return nil
}
