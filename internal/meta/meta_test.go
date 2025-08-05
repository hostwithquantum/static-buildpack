package meta_test

import (
	"os"
	"testing"

	"github.com/hostwithquantum/static-buildpack/internal/meta"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/stretchr/testify/assert"
)

func TestMeta(t *testing.T) {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel("debug")

	t.Run("Test web server", func(t *testing.T) {
		t.Setenv("BP_WEB_SERVER", "caddy")
		assert.Equal(t, "nginx", meta.DetectWebServer())
	})

	t.Run("detect go", func(t *testing.T) {
		assert.True(t, meta.NeedsGO("./../../tests/hugo-go", logEmitter))
		assert.False(t, meta.NeedsGO("./../../tests/hugo-example", logEmitter))
		assert.False(t, meta.NeedsGO("./../../tests/hugo-npm", logEmitter))
	})

	t.Run("detect npm", func(t *testing.T) {
		assert.False(t, meta.NeedsNPM("./../../tests/hugo-go", logEmitter))
		assert.False(t, meta.NeedsNPM("./../../tests/hugo-example", logEmitter))
		assert.True(t, meta.NeedsNPM("./../../tests/hugo-npm", logEmitter))
	})
}
