package detect_test

import (
	"os"
	"testing"

	"github.com/hostwithquantum/static-buildpack/internal/detect"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/stretchr/testify/assert"
)

func TestFinder(t *testing.T) {
	finder := detect.NewFinder(scribe.NewEmitter(os.Stdout).WithLevel("debug"))

	t.Run("Hugo", func(t *testing.T) {
		err := finder.Find("./../../tests/hugo-example")
		assert.NoError(t, err, "expected no error, got %s", err)
		assert.Equal(t, detect.HugoType, finder.GetStaticType())
	})

	t.Run("MdBook", func(t *testing.T) {
		err := finder.Find("./../../tests/mdbook-example")
		assert.NoError(t, err, "expected no error, got %s", err)
		assert.Equal(t, detect.MdBookType, finder.GetStaticType())
	})

	t.Run("NoMatch", func(t *testing.T) {
		err := finder.Find("./")
		assert.Error(t, err)
	})
}
