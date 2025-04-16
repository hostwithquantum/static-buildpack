package detect_test

import (
	"os"
	"testing"

	"github.com/hostwithquantum/static-buildpack/internal/detect"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/stretchr/testify/assert"
)

func TestFinder(t *testing.T) {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel("debug")

	finder := detect.NewFinder(logEmitter)

	t.Run("Hugo", func(t *testing.T) {
		err := finder.Find("./../../tests/hugo-example")
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		assert.Equal(t, detect.HugoType, finder.GetStaticType())
	})

	t.Run("MdBook", func(t *testing.T) {
		err := finder.Find("./../../tests/mdbook-example")
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}

		assert.Equal(t, detect.MdBookType, finder.GetStaticType())
	})

	t.Run("NoMatch", func(t *testing.T) {
		assert.Error(t, finder.Find("./"))
	})
}
