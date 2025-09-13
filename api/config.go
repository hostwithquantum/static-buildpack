package api

import (
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/packit/v2"
)

func decode(from io.Reader, to *BuildpackTOML) error {
	if _, err := toml.NewDecoder(from).Decode(to); err != nil {
		return packit.Fail.WithMessage("failed to decode buildpack.toml: %s", err)
	}
	return nil
}

// GetDefaultVersion retrieves the default version for a tool
// Priority: 1. Environment variable, 2. buildpack.toml metadata
// TODO(till): this is not very efficient, but maybe it doesn't matter
func GetDefault(cnbPath, envVar string) string {
	// First check environment variable
	if version := os.Getenv(envVar); version != "" {
		return version
	}

	// Fallback to buildpack.toml metadata
	bp, err := os.Open(filepath.Join(cnbPath, "buildpack.toml"))
	if err != nil {
		// If we can't read buildpack.toml, we have no fallback
		return ""
	}
	defer bp.Close()

	var config BuildpackTOML
	if err := decode(bp, &config); err != nil {
		return ""
	}

	// Find the environment variable and return its default
	for _, env := range config.Metadata.Configurations {
		if env.Name == envVar {
			return env.Default
		}
	}

	return ""
}
