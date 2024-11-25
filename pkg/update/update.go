package update

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/blang/semver/v4"
)

var (
	Version = ""
	Repo    = ""
)

func Check() error {
	available, url, err := check()
	if err != nil {
		fmt.Printf("Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	if available {
		fmt.Printf("An update is available: %s\n", url)
		fmt.Println("Press Enter to continue...")
		waitForEnter()
	}

	return nil
}

func check() (bool, string, error) {
	resp, err := http.Get(Repo)
	if err != nil {
		return false, "", fmt.Errorf("error checking for updates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("failed to fetch release information. Status: %v", resp.Status)
	}

	var release struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, "", fmt.Errorf("error parsing release information: %v", err)
	}

	Version = normalizeVersion(Version)
	latestVersion := normalizeVersion(release.TagName)

	current, err := semver.Parse(Version)
	if err != nil {
		return false, "", fmt.Errorf("invalid current version (%s): %v", Version, err)
	}

	latest, err := semver.Parse(latestVersion)
	if err != nil {
		return false, "", fmt.Errorf("invalid latest version (%s): %v", release.TagName, err)
	}

	if latest.GT(current) {
		return true, release.HTMLURL, nil
	}

	return false, "", nil
}

// normalizeVersion removes the leading 'v' if present
func normalizeVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		return strings.TrimPrefix(version, "v")
	}
	return version
}

func waitForEnter() {
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}
