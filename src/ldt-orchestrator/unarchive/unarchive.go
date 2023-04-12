package unarchive

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"strings"
)

func Untar(src, dest string) (string, error) {
	folder, err := prepareFolder(src, dest)
	if err != nil {
		log.Fatal("Unable to create folder", err)
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
		unpacked, err := os.Create(dest)
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
	subfolder := strings.Split(strings.Split(src, "/")[1], ".")[0]
	folder := dest + "/" + subfolder

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.Mkdir(folder, 0777); err != nil {
			return "", err
		}
	}
	return folder, nil
}
