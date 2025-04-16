package main

import (
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/fs"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

type (
	finder struct {
		matched bool
		matches map[string]string
		Files   []string
	}

	StaticType string
)

const (
	HugoType   StaticType = "hugo"
	MdBookType StaticType = "mdbook"
)

func Factory() *finder {
	finder := &finder{
		matched: false,
		matches: make(map[string]string),
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
	for _, metaFile := range f.Files {
		path := filepath.Join(workingDir, metaFile)
		if exist, err := fs.Exists(path); err != nil {
			return err
		} else if exist {
			f.matched = true
			f.matches[metaFile] = path
		}
	}
	return nil
}

func (f *finder) GetStaticType() StaticType {
	if _, ok := f.matches["book.toml"]; ok {
		return MdBookType
	}
	return HugoType
}

func (f *finder) HasMatch() bool {
	return f.matched
}

func detect(logs scribe.Emitter) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		logs.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		finder := Factory()
		if err := finder.Find(context.WorkingDir); err != nil {
			return packit.DetectResult{}, err
		}

		if !finder.HasMatch() {
			return packit.DetectResult{}, packit.Fail.WithMessage("no static site configuration found")
		}

		staticType := finder.GetStaticType()
		requirements := []packit.BuildPlanRequirement{
			{
				Name: "static",
				Metadata: map[string]interface{}{
					"type": staticType,
				},
			},
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Requires: requirements,
			},
		}, nil
	}
}

func main() {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))
	packit.Detect(detect(logEmitter))
}
