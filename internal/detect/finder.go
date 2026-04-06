package detect

import (
	"path/filepath"
	"slices"

	"github.com/paketo-buildpacks/packit/v2"
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
	if wdOK, err := fs.Exists(workingDir); err != nil {
		return err
	} else if !wdOK {
		return packit.Fail.WithMessage("WorkingDir does not exist: %s", workingDir)
	}

	f.log.Process("Detecting static site configuration")
	for _, metaFile := range f.Files {
		f.log.Subprocess(metaFile)

		// Build list of paths to check: root dir first, then config/_default/ for Hugo files
		paths := []string{filepath.Join(workingDir, metaFile)}
		if f.isHugoConfig(metaFile) {
			paths = append(paths, filepath.Join(workingDir, "config", "_default", metaFile))
		}

		for _, path := range paths {
			if exist, err := fs.Exists(path); err != nil {
				return err
			} else if exist {
				f.log.Detail("Found: %s", path)
				f.matched = true
				f.matches[metaFile] = path
				return nil
			}
		}
	}
	return packit.Fail.WithMessage("no static site configuration found in: %s", workingDir)
}

func (f *finder) isHugoConfig(filename string) bool {
	hugoConfigs := []string{"hugo.yml", "hugo.yaml", "hugo.toml", "hugo.json"}
	return slices.Contains(hugoConfigs, filename)
}

func (f *finder) GetStaticType() StaticType {
	if _, ok := f.matches["book.toml"]; ok {
		return MdBookType
	}
	return HugoType
}
