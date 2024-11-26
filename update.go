package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/blang/semver/v4"
)

var (
	Version = "" // This will be set by the linker
	Repo    = "" // This will be set by the linker
)

type UpdateInfo struct {
	Available bool   `json:"available"`
	Message   string `json:"message"`
	URL       string `json:"url"`
}

func CheckForUpdate() (UpdateInfo, error) {
	log.Printf("Checking for updates. Current version: %s, Repo: %s", Version, Repo)

	resp, err := http.Get(Repo)
	if err != nil {
		return UpdateInfo{Available: false, Message: fmt.Sprintf("Error checking for updates: %v", err)}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UpdateInfo{Available: false, Message: fmt.Sprintf("Failed to fetch release information. Status: %v", resp.Status)}, nil
	}

	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return UpdateInfo{Available: false, Message: fmt.Sprintf("Error parsing release information: %v", err)}, err
	}

	currentVersion := normalizeVersion(Version)
	latestVersion := normalizeVersion(release.TagName)

	current, err := semver.Parse(currentVersion)
	if err != nil {
		return UpdateInfo{Available: false, Message: fmt.Sprintf("Invalid current version (%s): %v", Version, err)}, err
	}

	latest, err := semver.Parse(latestVersion)
	if err != nil {
		return UpdateInfo{Available: false, Message: fmt.Sprintf("Invalid latest version (%s): %v", release.TagName, err)}, err
	}

	if latest.GT(current) {
		log.Printf("Update available: %s -> %s", currentVersion, latestVersion)

		return UpdateInfo{
			Available: true,
			Message:   release.TagName,
			URL:       release.HTMLURL,
		}, nil
	}

	log.Printf("You are using the latest version: %s", currentVersion)

	return UpdateInfo{
		Available: false,
		Message:   "You are using the latest version.",
	}, nil
}

// normalizeVersion removes the leading 'v' if present
func normalizeVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		return strings.TrimPrefix(version, "v")
	}
	return version
}
