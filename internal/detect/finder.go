package detect

import (
	"fmt"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

type (
	finder struct {
		matched bool
		matches map[string]string
		Files   []string
		log     scribe.Emitter
	}

	StaticType string
)

const (
	HugoType   StaticType = "hugo"
	MdBookType StaticType = "mdbook"
)

func NewFinder(log scribe.Emitter) *finder {
	finder := &finder{
		matched: false,
		matches: make(map[string]string),
		log:     log,
	}

	finder.Files = []string{
		// hugo
		"hugo.yml", "hugo.yaml", "hugo.toml", "hugo.json",
		// mdbook
		"book.toml",
	}

	return finder
}

func (f *finder) Find(workingDir string) error {
	f.log.Process("Detecting static site configuration")
	for _, metaFile := range f.Files {
		f.log.Subprocess(metaFile)
		path := filepath.Join(workingDir, metaFile)
		if exist, err := fs.Exists(path); err != nil {
			return err
		} else if exist {
			f.log.Detail("Found: %s", path)
			f.matched = true
			f.matches[metaFile] = path
			return nil
		}
	}
	return fmt.Errorf("no static site configuration found")
}

func (f *finder) GetStaticType() StaticType {
	if _, ok := f.matches["book.toml"]; ok {
		return MdBookType
	}
	return HugoType
}
