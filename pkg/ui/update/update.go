package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/blang/semver/v4"
	"github.com/gdamore/tcell/v2"

	"github.com/andrewrynhard-audio/bpm/pkg/state"
)

var (
	Version = ""
	Repo    = ""
)

type Update struct {
	sync.Once

	message string
}

func New() *Update {
	return &Update{
		message: Version,
	}
}

func (u *Update) Render(sharedState *state.State, screen tcell.Screen) {
	// Run the update check asynchronously with sync.Once
	u.Do(func() {
		go func() {
			available, url, err := check()
			if err != nil {
				u.message = fmt.Sprintf("Error checking for updates: %v", err)
				return
			}

			if available {
				u.message = fmt.Sprintf("An update is available: %s", url)
			}
		}()
	})

	// Display the update message
	renderText(screen, 1, 1, u.message, tcell.StyleDefault.Foreground(tcell.ColorYellow).Dim(true))
}

func (u *Update) Reset(sharedState *state.State, screen tcell.Screen) {
	// No-op
}

func (u *Update) StateChanged(sharedState *state.State, screen tcell.Screen) {
	// No-op
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

func renderText(screen tcell.Screen, x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		screen.SetContent(x+i, y, ch, nil, style)
	}
}

// normalizeVersion removes the leading 'v' if present
func normalizeVersion(version string) string {
	if strings.HasPrefix(version, "v") {
		return strings.TrimPrefix(version, "v")
	}
	return version
}
