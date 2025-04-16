package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
)

const (
	HugoLatestVersion   = "0.123.8"
	MdBookLatestVersion = "0.10.2"
)

func InstallHugo(layer packit.Layer, version string) error {
	if version == "latest" {
		version = HugoLatestVersion
	}

	// Create bin directory in the layer
	binDir := filepath.Join(layer.Path, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Download and install Hugo
	hugoPath := filepath.Join(binDir, "hugo")
	if err := downloadAndInstall(getHugoURL(version), hugoPath); err != nil {
		return fmt.Errorf("failed to download Hugo: %w", err)
	}

	return nil
}

func InstallMdBook(layer packit.Layer, version string) error {
	if version == "latest" {
		version = MdBookLatestVersion
	}

	// Create bin directory in the layer
	binDir := filepath.Join(layer.Path, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	// Download and install mdBook
	mdbookPath := filepath.Join(binDir, "mdbook")
	if err := downloadAndInstall(getMdBookURL(version), mdbookPath); err != nil {
		return fmt.Errorf("failed to download mdBook: %w", err)
	}

	return nil
}

func GetStaticToolPath(layer packit.Layer, tool string) string {
	return filepath.Join(layer.Path, "bin", tool)
}
