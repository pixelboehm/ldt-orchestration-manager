package unarchive

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Untar(t *testing.T) {
	require := require.New(t)

	source, err := download_helper("https://github.com/pixelboehm/ldt/releases/download/v0.6.0/lightbulb_Darwin_x86_64.tar.gz")
	require.NoError(err)
	var dest string = "./out/test_ldt"

	_, err = Untar(source, dest)
	require.NoError(err)

	os.RemoveAll("./out")
	os.RemoveAll("./resources")
}

func download_helper(address string) (string, error) {
	url, _ := url.Parse(address)
	filename := strings.Split(url.Path, "/")[6]
	file, err := create("./resources/" + filename)
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
