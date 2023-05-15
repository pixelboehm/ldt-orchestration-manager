package unarchive

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func Untar(src, dest string) (string, error) {
	folder, err := prepareFolder(src, dest)
	if err != nil {
		panic(fmt.Sprintf("Unable to create folder %v\n", err))
	}

	file, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzip, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzip.Close()

	tar := tar.NewReader(gzip)

	var files []string

	for {
		nextFile, err := tar.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("failed to read next tar entry")
			return "", err
		}
		dest := folder + "/" + nextFile.Name
		files = append(files, dest)
		unpacked, err := create(dest)
		if err != nil {
			log.Printf("failed to create: %s\n", dest)
			return "", err
		}
		defer unpacked.Close()

		_, err = io.Copy(unpacked, tar)
		if err != nil {
			log.Printf("failed to unpack to: %s\n", unpacked.Name())
			return "", err
		}
	}
	return files[1], nil
}

func prepareFolder(src, dest string) (string, error) {
	folder := dest

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0777); err != nil {
			return "", err
		}
	}
	return folder, nil
}

func create(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return nil, err
	}
	return os.Create(path)
}
