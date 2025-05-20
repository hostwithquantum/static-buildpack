package api

import (
	"os"
	"path/filepath"
)

// GetWorkingDir determines the complete path of where the static site is located.
// By default, this will be the "WorkingDir" from build/detect Context, but it may
// be extended via config.
func GetWorkingDir(defaultWorkDir string) string {
	if path := os.Getenv(StaticPathEnv); path != "" {
		return filepath.Join(defaultWorkDir, path)
	}
	return defaultWorkDir
}
