package ldtorchestrator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	setupDiscoveryConfig()
	code := m.Run()
	teardownDiscoveryConfig()
	os.Exit(code)
}

var c *DiscoveryConfig = nil

func setupDiscoveryConfig() {
	d1 := []byte("https://github.com/pixelboehm/ldt\n#https://github.com/pixelboehm/longevity\n")
	if err := os.WriteFile("./config", d1, 0644); err != nil {
		panic(err)
	}
	c = NewConfig("./config")
}

func teardownDiscoveryConfig() {
	os.Remove("./config")
}

func ensureConfigExists(t *testing.T) {
	require := require.New(t)
	if c != nil {
		require.FileExists(c.repository_file)
	} else {
		require.Fail("DiscoveryConfig is not initialized")
	}
}

func Test_UpdateRepositories(t *testing.T) {
	ensureConfigExists(t)
	assert := assert.New(t)

	expected := 1
	c.updateRepositories()
	actual := len(c.repositories)
	assert.Equal(expected, actual)
}

func Test_IsGithubRepository(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		desc  string
		input string
		want  bool
	}{
		{
			desc:  "github url with HTTPS",
			input: "https://github.com/foobar",
			want:  true,
		}, {
			desc:  "github url with HTTPS and www",
			input: "https://www.github.com/foobar",
			want:  true,
		}, {
			desc:  "github url without HTTPS but with www",
			input: "www.github.com/foobar",
			want:  true,
		}, {
			desc:  "github url without HTTPS and without www",
			input: "github.com/foobar",
			want:  true,
		}, {
			desc:  "some URL that is not related to github",
			input: "https://www.google.com/foobar",
			want:  false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := isGithubRepository(tC.input)
			assert.Equal(tC.want, got)
		})
	}
}

func Test_DiscoverLDTs(t *testing.T) {
	ensureConfigExists(t)
	assert := assert.New(t)

	c.DiscoverLDTs()
	assert.NotNil(len(c.supportedLDTs.LDTs))
}
