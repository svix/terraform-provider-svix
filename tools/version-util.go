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
	if len(os.Args) < 2 {
		log.Fatal("Chose subcommand: `bump` or `check`")
	}
	subCommand := os.Args[1]
	switch subCommand {
	case "bump":
		bumpVersion()
	case "check":
		checkVersion()
	default:
		log.Fatal("Chose subcommand: `bump` or `check`")
	}

}

// this will run in CI right before we publish a release
func checkVersion() {
	currentVersion, err := getCurrentVersion()
	if err != nil {
		log.Fatal(err)
	}
	versionFromGitTag, err := getVersionFromCliArgs()
	if err != nil {
		log.Fatal(err)
	}

	// this will strip the leading `v`
	if currentVersion.String() != versionFromGitTag.String() {
		log.Fatalf("Version from git tag `%s` does not equal current version `%s`", versionFromGitTag, currentVersion)
	}

	for _, path := range filesToReplace {
		err = checkFile(path, currentVersion.String())
		if err != nil {
			log.Fatal(err)
		}
	}

}

func checkFile(path string, currentVersion string) error {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(fileContent), currentVersion) {
		log.Fatalf("File `%s` is missing expected version `%s`", path, currentVersion)
	}
	return nil
}

func bumpVersion() {
	currentVersion, err := getCurrentVersion()
	if err != nil {
		log.Fatal(err)
	}
	newVersion, err := getVersionFromCliArgs()
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

func getVersionFromCliArgs() (*semver.Version, error) {
	if len(os.Args) < 3 {
		log.Fatal("Usage: version bump <new version>")
	}

	newVersionStr := os.Args[2]
	return semver.NewVersion(strings.TrimSpace(string(newVersionStr)))

}
