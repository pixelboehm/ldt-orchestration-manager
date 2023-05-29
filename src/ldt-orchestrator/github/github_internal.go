package github

import (
	"crypto/sha256"
	"fmt"
	. "longevity/src/types"
	"net/url"
	"strings"
)

func filterLDTInformationFromURL(address string) *LDT {
	u, _ := url.Parse(address)
	user := strings.Split(u.Path, "/")[1]

	version := strings.Split(u.Path, "/")[5]

	filename := strings.Split(u.Path, "/")[6]
	withoutSuffix := strings.Split(filename, ".")[0]

	ldtname, rest, _ := strings.Cut(withoutSuffix, "_")
	os, arch, _ := strings.Cut(rest, "_")

	switch arch {
	case "x86_64":
		arch = "amd64"
	}

	ldt := &LDT{
		Name:    ldtname,
		User:    user,
		Version: version,
		Os:      strings.ToLower(os),
		Arch:    arch,
		Url:     address,
	}
	return ldt
}

func finalizeLDT(ldt *LDT) {
	ldt.Hash = createHash(ldt)
}

func createHash(l *LDT) []byte {
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", l)))
	return h.Sum(nil)
}

func isArchive(file string) bool {
	return strings.HasSuffix(file, ".tar.gz")
}

func parseRepository(repo string) (string, string) {
	split := strings.Split(repo, "/")
	return split[3], split[4]
}
