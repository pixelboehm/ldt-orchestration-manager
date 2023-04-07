package github

import (
	"testing"

	. "longevity/src/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CreateGithubDiscovererWithoutToken(t *testing.T) {
	assert := assert.New(t)
	token := "NOT_EXISTING_TOKEN"
	gd := NewGithubClient(token)
	assert.False(gd.Authenticated)
}

func TestCreateGithubDiscovererWithToken(t *testing.T) {
	assert := assert.New(t)
	token := "GH_ACCESS_TOKEN"
	gd := NewGithubClient(token)
	assert.True(gd.Authenticated)
}

func Test_FilterLDTsFromReleases(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	gd := NewGithubClient("GH_ACCESS_TOKEN")
	require.True(gd.Authenticated)
	releases := gd.GetReleasesFromRepository("pixelboehm", "ldt")
	assert.NotEmpty(releases)

	ldt_list := gd.FilterLDTsFromReleases(releases)
	assert.NotEmpty(ldt_list.LDTs)
}

func Test_FilteringLDTInformationFromURL(t *testing.T) {
	assert := assert.New(t)

	var testCases = []struct {
		name  string
		input string
		want  LDT
	}{
		{".tar.gz with x86_64 architecture", "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.tar.gz", LDT{Name: "pixelboehm/switch", Version: "v0.2.1", Os: "Darwin", Arch: "x86_64", Url: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.tar.gz"}},

		{".tar.gz with arm64 architecture", "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.tar.gz", LDT{Name: "pixelboehm/switch", Version: "v0.2.1", Os: "Darwin", Arch: "arm64", Url: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.tar.gz"}},

		{".zip with x86_64 architecture", "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.zip", LDT{Name: "pixelboehm/switch", Version: "v0.2.1", Os: "Darwin", Arch: "x86_64", Url: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.zip"}},

		{".zip with arm64 architecture", "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.zip", LDT{Name: "pixelboehm/switch", Version: "v0.2.1", Os: "Darwin", Arch: "arm64", Url: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.zip"}},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := filterLDTInformationFromURL(tC.input)
			assert.Equal(tC.want, got)
		})
	}
}
