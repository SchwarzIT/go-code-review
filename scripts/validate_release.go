package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
)

const releaseFile = "./release.toml"

type Release struct {
	Metadata struct {
		Author             string `toml:"Author"`
		ReleaseFileVersion string `toml:"Release_file_version"`
	} `toml:"Metadata"`

	Description struct {
		Notes string `toml:"Notes"`
	} `toml:"Description"`

	Digest struct {
		Features     []Section `toml:"Features"`
		Improvements []Section `toml:"Improvements"`
		Bugs         []Section `toml:"Bugs"`
	} `toml:"Digest"`
}

type Section struct {
	Name        string `toml:"Name"`
	Issue       int    `toml:"Issue"`
	Description string `toml:"Description"`
}

func validateRelease() error {
	content, err := os.ReadFile(releaseFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	if !strings.HasPrefix(string(content), "[Metadata]") {
		return fmt.Errorf("file must start with [Metadata] section")
	}

	var release Release
	err = toml.Unmarshal(content, &release)
	if err != nil {
		return fmt.Errorf("invalid TOML format: %v", err)
	}

	if release.Metadata.Author == "" || release.Metadata.ReleaseFileVersion == "" {
		return fmt.Errorf("missing required fields in Metadata")
	}

	if release.Description.Notes == "" {
		return fmt.Errorf("Notes in Description must not be empty")
	}

	if len(release.Digest.Features) == 0 && len(release.Digest.Improvements) == 0 && len(release.Digest.Bugs) == 0 {
		return fmt.Errorf("at least one Digest section (Features, Improvements, Bugs) must be provided")
	}

	validateSection := func(sections []Section) error {
		for _, section := range sections {
			if section.Name == "" || section.Issue == 0 || section.Description == "" {
				return fmt.Errorf("Digest entries must have Name, Issue, and Description filled out")
			}
		}
		return nil
	}

	if err := validateSection(release.Digest.Features); err != nil {
		return err
	}
	if err := validateSection(release.Digest.Improvements); err != nil {
		return err
	}
	if err := validateSection(release.Digest.Bugs); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Checking if release file exists")
	if _, err := os.Stat(releaseFile); os.IsNotExist(err) {
		fmt.Println("release.toml file not found")
		os.Exit(1)
	}

	fmt.Println("Validating release file")
	err := validateRelease()
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		os.Exit(1)
	}

}
