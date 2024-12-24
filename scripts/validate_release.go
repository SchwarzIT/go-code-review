package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
)

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

func validateRelease(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Basic format check
	if !strings.HasPrefix(string(content), "[Metadata]") {
		return fmt.Errorf("file must start with [Metadata] section")
	}

	// Parse TOML
	var release Release
	err = toml.Unmarshal(content, &release)
	if err != nil {
		return fmt.Errorf("invalid TOML format: %v", err)
	}

	// Check required fields in Metadata
	if release.Metadata.Author == "" || release.Metadata.ReleaseFileVersion == "" {
		return fmt.Errorf("missing required fields in Metadata")
	}

	// Check Notes in Description
	if release.Description.Notes == "" {
		return fmt.Errorf("Notes in Description must not be empty")
	}

	// Check Digest sections
	if len(release.Digest.Features) == 0 && len(release.Digest.Improvements) == 0 && len(release.Digest.Bugs) == 0 {
		return fmt.Errorf("at least one Digest section (Features, Improvements, Bugs) must be provided")
	}

	// Ensure every entry in Digest has required fields
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
	releaseFile := "./release.toml"
	// Check if file exists

	fmt.Println("Checking if release file exists")
	if _, err := os.Stat(releaseFile); os.IsNotExist(err) {
		fmt.Println("release.toml file not found")
		os.Exit(1)
	}

	// Validate file
	fmt.Println("Validating release file")
	err := validateRelease(releaseFile)
	if err != nil {
		fmt.Printf("Validation failed: %v\n", err)
		os.Exit(1)
	}

}
