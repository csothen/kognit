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
	"strings"
)

type ItemType int
type DirectoryCompressionAlgorithm int
type FileCompressionAlgorithm int
type ImageCompressionAlgorithm int

type CompressionAlgorithm interface {
	Encode(src string) error
	Decode(src string) error
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

func (a DirectoryCompressionAlgorithm) Encode(src string) error {
	dest := src

	switch a {
	case ZIP:
		dest += ".zip"
		if err := encodeZipArchive(src, dest); err != nil {
			return err
		}
	case TAR:
		dest += ".tar.gz"
		if err := encodeTarArchive(src, dest); err != nil {
			return err
		}
	}
	return nil
}

func encodeZipArchive(src, dest string) error {
	files, err := allDirFiles(src)
	if err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	zipWriter := zip.NewWriter(f)
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

func encodeTarArchive(src, dest string) error {
	files, err := allDirFiles(src)
	if err != nil {
		return nil
	}

	tarFile, err := os.Create(dest)
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

func allDirFiles(src string) ([]string, error) {
	files := []string{}

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
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

func (a DirectoryCompressionAlgorithm) Decode(src string) error {
	dest := filepath.Dir(src)
	switch a {
	case ZIP:
		if err := decodeZipArchive(src, dest); err != nil {
			return err
		}
	case TAR:
		if err := decodeTarArchive(src, dest); err != nil {
			return err
		}
	}
	return nil
}

func decodeZipArchive(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	os.MkdirAll(dest, 0755)

	for _, f := range r.File {
		err := extractFromZip(f, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFromZip(f *zip.File, dest string) error {
	file, err := f.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	path := filepath.Join(dest, f.Name)

	if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
		return fmt.Errorf("illegal file path: %s", path)
	}

	if f.FileInfo().IsDir() {
		os.MkdirAll(path, f.Mode())
	} else {
		os.MkdirAll(filepath.Dir(path), f.Mode())
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func decodeTarArchive(src, dest string) error {
	stream, err := os.Open(src)
	if err != nil {
		return err
	}
	defer stream.Close()

	gzipReader, err := gzip.NewReader(stream)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	r := tar.NewReader(gzipReader)

	for true {
		header, err := r.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		err = extractFromTar(r, header, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

func extractFromTar(r *tar.Reader, header *tar.Header, dest string) error {
	switch header.Typeflag {
	case tar.TypeDir:
		if err := os.Mkdir(header.Name, 0755); err != nil {
			return err
		}
	case tar.TypeReg:
		outFile, err := os.Create(header.Name)
		if err != nil {
			return err
		}

		if _, err := io.Copy(outFile, r); err != nil {
			return err
		}
		outFile.Close()

	default:
		return errors.New("Unknown header type")
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
