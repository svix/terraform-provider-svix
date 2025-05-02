package main

import (
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// relative to the location of this file
var filesToReplace = []string{
	"../examples/provider/provider.tf",
	"../internal/version.go",
	"../.version",
}

func main() {
	currentVersion, err := getCurrentVersion()
	if err != nil {
		log.Fatal(err)
	}
	newVersion, err := getNewVersion()
	if err != nil {
		log.Fatal(err)
	}

	if newVersion.Compare(currentVersion) == 0 {
		log.Fatalf("New version `%s` must not be equal to current version `%s`", newVersion, currentVersion)
	}
	if newVersion.Compare(currentVersion) == -1 {
		log.Fatalf("New version `%s` must be greater then the current version `%s`", newVersion, currentVersion)
	}
	for _, path := range filesToReplace {
		err = replaceFile(path, currentVersion.String(), newVersion.String())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func replaceFile(path string, oldVersion string, newVersion string) error {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	updatedFile := strings.Replace(string(fileContent), oldVersion, newVersion, 1)
	return os.WriteFile(path, []byte(updatedFile), fs.FileMode(0644))
}

func getCurrentVersion() (*semver.Version, error) {
	content, err := os.ReadFile("../.version")
	if err != nil {
		return nil, err
	}
	return semver.NewVersion(strings.TrimSpace(string(content)))
}

func getNewVersion() (*semver.Version, error) {
	if len(os.Args) < 2 {
		log.Fatal("Usage: bump-version <new version>")
	}

	newVersionStr := os.Args[1]
	return semver.NewVersion(strings.TrimSpace(string(newVersionStr)))

}
