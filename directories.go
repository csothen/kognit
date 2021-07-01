package kognit

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
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
		if info.Mode().IsRegular() {
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

	if f.FileInfo().IsDir() {
		os.MkdirAll(path, 0755)
	} else {
		os.MkdirAll(filepath.Dir(path), 0755)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
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

	os.MkdirAll(dest, 0755)

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
		os.MkdirAll(filepath.Dir(header.Name), 0755)
		f, err := os.OpenFile(header.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(f, r); err != nil {
			return err
		}

	default:
		return errors.New("Unknown header type")
	}
	return nil
}
