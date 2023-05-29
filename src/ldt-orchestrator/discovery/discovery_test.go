package discovery

import (
	"longevity/src/types"
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

var c *Discoverer = nil

func setupDiscoveryConfig() {
	d1 := []byte("https://github.com/pixelboehm/ldt\n#https://github.com/pixelboehm/longevity\n")
	if err := os.WriteFile("./config", d1, 0644); err != nil {
		panic(err)
	}
	c = NewDiscoverer("./config")
}

func teardownDiscoveryConfig() {
	os.Remove("./config")
}

func ensureConfigExists(t *testing.T) {
	require := require.New(t)
	if c != nil {
		require.FileExists(c.repository_source)
	} else {
		require.Fail("DiscoveryConfig is not initialized")
	}
}

func Test_UpdateRepositories(t *testing.T) {
	ensureConfigExists(t)
	assert := assert.New(t)
	require := require.New(t)

	expected := 1
	err := c.updateRepositories()
	actual := len(c.repositories)
	require.NoError(err)
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
	t.Skip("skipping test")
	ensureConfigExists(t)
	assert := assert.New(t)

	c.DiscoverLDTs()
	assert.NotNil(len(c.SupportedLDTs.LDTs))
}

func Test_updateLatestTag(t *testing.T) {
	ensureConfigExists(t)
	assert := assert.New(t)
	require := require.New(t)

	injectLDTFromFakeUserIn(&c.SupportedLDTs.LDTs)
	c.DiscoverLDTs()
	require.NotNil(len(c.SupportedLDTs.LDTs))

	var actual_latest_tags int = 0
	for _, ldt := range c.SupportedLDTs.LDTs {
		if ldt.Version == "latest" {
			actual_latest_tags += 1
		}
	}
	require.NotZero(actual_latest_tags)

	var expected_latest_tags int = getUniqueUserLDTCombinations(c.SupportedLDTs.LDTs)
	assert.Equal(expected_latest_tags, actual_latest_tags)
}

func injectLDTFromFakeUserIn(list *[]types.LDT) {
	ldt := &types.LDT{
		Name:    "switch",
		User:    "fake_user",
		Version: "v0.10.2",
		Os:      "darwin",
		Arch:    "amd64",
		Url:     "",
		Hash:    nil,
	}

	*list = append([]types.LDT{*ldt}, *list...)
}

func getUniqueUserLDTCombinations(ldts []types.LDT) int {
	var unique []types.LDT
loop:
	for _, l := range ldts {
		for i, u := range unique {
			if l.Name == u.Name && l.User == u.User {
				unique[i] = l
				continue loop
			}
		}
		unique = append(unique, l)
	}
	return len(unique)
}

func Test_isURL(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		desc  string
		input string
		want  bool
	}{
		{
			desc:  "is a github URL",
			input: "https://github.com/pixelboehm/ldt",
			want:  true,
		},
		{
			desc:  "is a filepath",
			input: "$HOME/.local/etc/orchstration-manager/repositories.list",
			want:  false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := isURL(tC.input)
			assert.Equal(tC.want, res)
		})
	}
}

// Note: This test has is dirty and needs to be done better. Only comparing the first 59 characters, because os and arch change depending on the system
func Test_GetURLFromLDTByName(t *testing.T) {
	ensureConfigExists(t)

	assert := assert.New(t)
	require := require.New(t)

	user := "pixelboehm"
	ldt := "lightbulb"
	version := "v0.5.0"
	want := "https://github.com/pixelboehm/ldt/releases/download/v0.5.0/lightbulb_Darwin_x86_64.tar.gz"
	got, err := c.GetURLFromLDTByName(user, ldt, version)
	require.NoError(err)

	assert.Equal(want[:59], got[:59])
}
