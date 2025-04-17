package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

const (
	HugoLatestVersion   = "0.123.8"
	MdBookLatestVersion = "0.10.2"
)

func installHugo(log scribe.Emitter, layer packit.Layer, version string) error {
	if version == "latest" {
		version = HugoLatestVersion
	}

	// Create bin directory in the layer
	binDir := filepath.Join(layer.Path, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Download and install Hugo
	if err := downloadAndInstall(log, getHugoURL(version), binDir); err != nil {
		return fmt.Errorf("failed to download Hugo: %w", err)
	}

	log.Subprocess("Installed hugo")

	return nil
}

func installMdBook(log scribe.Emitter, layer packit.Layer, version string) error {
	if version == "latest" {
		version = MdBookLatestVersion
	}

	// Create bin directory in the layer
	binDir := filepath.Join(layer.Path, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Download and install mdBook
	if err := downloadAndInstall(log, getMdBookURL(version), binDir); err != nil {
		return fmt.Errorf("failed to download mdBook: %w", err)
	}

	log.Subprocess("Installed mdbook")

	return nil
}
