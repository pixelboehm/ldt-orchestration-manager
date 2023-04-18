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
		desc  string
		input string
		want  *LDT
	}{
		{
			desc:  ".tar.gz with x86_64 (amd) architecture",
			input: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.tar.gz",
			want: &LDT{
				Name:    "pixelboehm/switch",
				Version: "v0.2.1",
				Os:      "darwin",
				Arch:    "amd64",
				Url:     "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.tar.gz",
				Hash:    nil},
		},
		{
			desc:  ".tar.gz with arm64 architecture",
			input: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.tar.gz",
			want: &LDT{
				Name:    "pixelboehm/switch",
				Version: "v0.2.1",
				Os:      "darwin",
				Arch:    "arm64",
				Url:     "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.tar.gz",
				Hash:    nil},
		},
		{
			desc:  ".zip with x86_64 (amd64) architecture",
			input: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.zip",
			want: &LDT{
				Name:    "pixelboehm/switch",
				Version: "v0.2.1",
				Os:      "darwin",
				Arch:    "amd64",
				Url:     "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.zip",
				Hash:    nil},
		},
		{
			desc:  ".zip with arm64 architecture",
			input: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.zip",
			want: &LDT{
				Name:    "pixelboehm/switch",
				Version: "v0.2.1",
				Os:      "darwin",
				Arch:    "arm64",
				Url:     "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.zip",
				Hash:    nil},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := filterLDTInformationFromURL(tC.input)
			assert.Equal(tC.want, got)
		})
	}
}
