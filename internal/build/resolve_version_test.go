package build_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/hostwithquantum/static-buildpack/api"
	"github.com/hostwithquantum/static-buildpack/internal/build"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/stretchr/testify/assert"
)

func TestResolveVersion(t *testing.T) {
	cnbPath := "../../" // project root containing buildpack.toml

	t.Run("env var set returns env var value", func(t *testing.T) {
		t.Setenv(api.HugoVersionEnv, "0.999.0")

		layer := packit.Layer{
			Metadata: map[string]interface{}{"version": "0.111.0"},
		}
		log := scribe.NewEmitter(os.Stdout)

		got := build.ResolveVersion(log, layer, cnbPath, api.HugoVersionEnv)
		assert.Equal(t, "0.999.0", got)
	})

	t.Run("no env var with cached version returns cached", func(t *testing.T) {
		var buf bytes.Buffer
		log := scribe.NewEmitter(&buf)

		layer := packit.Layer{
			Metadata: map[string]interface{}{"version": "0.155.0"},
		}

		got := build.ResolveVersion(log, layer, cnbPath, api.HugoVersionEnv)
		assert.Equal(t, "0.155.0", got)
		assert.Contains(t, buf.String(), "Reusing previously installed version 0.155.0")
	})

	t.Run("no env var no cache returns buildpack.toml default", func(t *testing.T) {
		layer := packit.Layer{
			Metadata: map[string]interface{}{},
		}
		log := scribe.NewEmitter(os.Stdout)

		got := build.ResolveVersion(log, layer, cnbPath, api.HugoVersionEnv)
		assert.Equal(t, "0.160.1", got)
	})

	t.Run("nil metadata returns buildpack.toml default", func(t *testing.T) {
		layer := packit.Layer{}
		log := scribe.NewEmitter(os.Stdout)

		got := build.ResolveVersion(log, layer, cnbPath, api.HugoVersionEnv)
		assert.Equal(t, "0.160.1", got)
	})

	t.Run("empty cached version string returns buildpack.toml default", func(t *testing.T) {
		layer := packit.Layer{
			Metadata: map[string]interface{}{"version": ""},
		}
		log := scribe.NewEmitter(os.Stdout)

		got := build.ResolveVersion(log, layer, cnbPath, api.HugoVersionEnv)
		assert.Equal(t, "0.160.1", got)
	})
}
