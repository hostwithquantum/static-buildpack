package api_test

import (
	"testing"

	"github.com/hostwithquantum/static-buildpack/api"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetDefault(t *testing.T) {
	ctx := packit.BuildContext{
		CNBPath: "../", // Point to the root directory for testing
	}

	t.Run("Test reading hugo version from buildpack.toml", func(t *testing.T) {
		version := api.GetDefault(ctx.CNBPath, "BP_RUNWAY_STATIC_HUGO_VERSION")
		assert.NotEmpty(t, version, "Expected to get hugo version from buildpack.toml, got empty string")
		assert.Equal(t, "0.147.4", version, "Expected hugo version 0.147.4, got %s", version)
	})

	t.Run("Test reading mdbook version from buildpack.toml", func(t *testing.T) {
		version := api.GetDefault(ctx.CNBPath, "BP_RUNWAY_STATIC_MDBOOK_VERSION")
		assert.NotEmpty(t, version, "Expected to get mdbook version from buildpack.toml, got empty string")
		assert.Equal(t, "0.4.49", version, "Expected mdbook version 0.4.49, got %s", version)
	})

	t.Run("Test with non-existent tool", func(t *testing.T) {
		version := api.GetDefault(ctx.CNBPath, "BP_NONEXISTENT")
		assert.Empty(t, version, "Expected empty string for nonexistent tool, got %s", version)
	})
}

func TestGetDefaultWithEnvVar(t *testing.T) {
	ctx := packit.BuildContext{
		CNBPath: "../",
	}

	// Test environment variable takes precedence
	t.Setenv("BP_RUNWAY_STATIC_HUGO_VERSION", "1.2.3")

	version := api.GetDefault(ctx.CNBPath, "BP_RUNWAY_STATIC_HUGO_VERSION")
	assert.Equal(t, "1.2.3", version, "Expected env var version 1.2.3, got %s", version)
}
