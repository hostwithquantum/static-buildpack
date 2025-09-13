package api

import (
	"path/filepath"
)

// GetWorkingDir determines the complete path of where the static site is located.
// By default, this will be the "WorkingDir" from build/detect Context, but it may
// be extended via config.
func GetWorkingDir(cnbPath, workDir string) string {
	path := GetDefault(cnbPath, StaticPathEnv)
	if path != "" {
		return filepath.Join(workDir, path)
	}

	return workDir
}
