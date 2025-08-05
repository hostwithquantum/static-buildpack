package meta

import (
	"os"
	"path/filepath"
	"slices"

	"github.com/paketo-buildpacks/packit/v2/scribe"
)

// Determine which web server to use
func DetectWebServer() string {
	webServer := os.Getenv("BP_WEB_SERVER")
	if slices.Contains([]string{"httpd", "nginx"}, webServer) {
		return webServer
	}
	return "nginx"
}

// some Hugo sites use go modules to manage themes
func NeedsGO(workingDir string, logs scribe.Emitter) bool {
	logs.Process("Check for go.mod")

	if _, err := os.Stat(filepath.Join(workingDir, "go.mod")); err == nil {
		logs.Detail("Found go.mod: go needed")
		return true
	}

	return false
}

// support node while building
func NeedsNPM(workingDir string, logs scribe.Emitter) bool {
	logs.Process("Check for node/npm")

	// this is either set by the Runway builder, or the user of this
	// buildpack has requested it themselves; in case you are not
	// yet a customer Runway — Hello to you!
	if _, ok := os.LookupEnv("BP_NODE_RUN_SCRIPTS"); ok {
		logs.Detail("Found BP_NODE_RUN_SCRIPTS: npm needed")
		return true
	}

	// check if we have a BP_NODE_PROJECT_PATH (in case)
	if path, ok := os.LookupEnv("BP_NODE_PROJECT_PATH"); ok {
		workingDir = filepath.Join(workingDir, path)
	}

	// check if the workingDir contains a package.json
	if _, err := os.Stat(filepath.Join(workingDir, "package.json")); err == nil {
		logs.Detail("Found package.json: npm needed")
		return true
	}

	logs.Detail("No npm")
	return false
}
