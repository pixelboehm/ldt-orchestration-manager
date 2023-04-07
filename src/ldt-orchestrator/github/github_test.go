package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetReleasesFromRepository(t *testing.T) {
	assert := assert.New(t)
	releases, err := GetReleasesFromRepository("pixelboehm", "ldt")
	assert.NoError(err)
	assert.NotEmpty(releases)
}

func Test_FilteringReleases(t *testing.T) {
	assert := assert.New(t)
	releases, _ := GetReleasesFromRepository("pixelboehm", "ldt")
	release := releases[0]
	res := filterLDTs(release)

	assert.True(res)
}

func Test_FilteringURL(t *testing.T) {
	assert := assert.New(t)
	url := "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.tar.gz"

	type Expected struct {
		name    string
		version string
		os      string
		arch    string
	}
	expected := Expected{name: "pixelboehm/switch", version: "v0.2.1", os: "Darwin", arch: "x86_64"}

	name, version, os, arch := filterURL(url)
	assert.Equal(expected.version, version)
	assert.Equal(expected.os, os)
	assert.Equal(expected.arch, arch)
	assert.Equal(expected.name, name)
}
