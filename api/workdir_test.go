package api_test

import (
	"path/filepath"
	"testing"

	"github.com/hostwithquantum/static-buildpack/api"
	"github.com/stretchr/testify/assert"
)

func TestGetWorkingDir(t *testing.T) {
	workDir := "/build/workspace"

	t.Run("default", func(t *testing.T) {
		assert.Equal(t, workDir, api.GetWorkingDir("../", workDir))
	})

	t.Run("env", func(t *testing.T) {
		t.Setenv(api.StaticPathEnv, "sub-dir")
		path := api.GetWorkingDir("../", workDir)
		assert.Equal(t, filepath.Join(workDir, "sub-dir"), path)
	})
}
