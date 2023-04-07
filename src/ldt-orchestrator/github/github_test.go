package github

import (
	"testing"

	. "longevity/src/ldt-orchestrator"

	"github.com/stretchr/testify/assert"
)

var gd = NewGithubDiscoverer("GH_ACCESS_TOKEN")

func Test_CreateGithubDiscovererWithoutToken(t *testing.T) {
	assert := assert.New(t)
	token := "NOT_EXISTING_TOKEN"
	gd := NewGithubDiscoverer(token)
	assert.False(gd.autenticated, false)
}

func TestCreateGithubDiscovererWithToken(t *testing.T) {
	assert := assert.New(t)
	token := "GH_ACCESS_TOKEN"
	gd := NewGithubDiscoverer(token)
	assert.False(gd.autenticated, true)
}

func Test_GetReleasesFromRepository(t *testing.T) {
	t.Skip()
	assert := assert.New(t)
	releases := gd.GetReleasesFromRepository("pixelboehm", "ldt")
	assert.NotEmpty(releases)
}

func Test_FilteringReleases(t *testing.T) {
	t.Skip("not implemented yet")
	// assert := assert.New(t)
	// releases := gd.GetReleasesFromRepository("pixelboehm", "ldt")
	// release := releases[0]
	// res := filterLDTs(release)

	// assert.True(res)
}

func Test_FilteringURL(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		name  string
		input string
		want  LDT
	}{
		{".tar.gz with x86_64 architecture", "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.tar.gz", LDT{Name: "pixelboehm/switch", Version: "v0.2.1", Os: "Darwin", Arch: "x86_64", Url: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.tar.gz"}},

		{".tar.gz with arm64 architecture", "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.tar.gz", LDT{Name: "pixelboehm/switch", Version: "v0.2.1", Os: "Darwin", Arch: "arm64", Url: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.tar.gz"}},

		{".zip with x86_64 architecture", "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.zip", LDT{Name: "pixelboehm/switch", Version: "v0.2.1", Os: "Darwin", Arch: "x86_64", Url: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_x86_64.zip"}},

		{".zip with arm64 architecture", "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.zip", LDT{Name: "pixelboehm/switch", Version: "v0.2.1", Os: "Darwin", Arch: "arm64", Url: "https://github.com/pixelboehm/ldt/releases/download/v0.2.1/switch_Darwin_arm64.zip"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterLDTInformationFromURL(tt.input)
			assert.Equal(tt.want, got)
		})
	}
}
