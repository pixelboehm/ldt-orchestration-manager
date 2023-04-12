package manager

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func download(address string) (string, error) {
	url, _ := url.Parse(address)
	filename := strings.Split(url.Path, "/")[6]

	file, err := os.Create("resources/" + filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	response, err := http.Get(address)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	name := file.Name()

	return name, nil
}
