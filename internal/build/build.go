package build

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hostwithquantum/static-buildpack/api"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/pexec"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

type (
	StaticType string
)

const (
	HugoType   StaticType = "hugo"
	MdBookType StaticType = "mdbook"

	// defaults
	HugoVersion   = "0.147.4"
	MdBookVersion = "0.4.49"
)

func Build(log scribe.Emitter) packit.BuildFunc {
	return func(ctx packit.BuildContext) (packit.BuildResult, error) {
		log.Title("%s %s", ctx.BuildpackInfo.Name, ctx.BuildpackInfo.Version)

		var plan packit.BuildpackPlanEntry
		for _, p := range ctx.Plan.Entries {
			if p.Name != "static-buildpack" {
				continue
			}

			plan = p
			break
		}

		rawType, ok := plan.Metadata["static-type"].(string)
		if !ok {
			return packit.BuildResult{}, fmt.Errorf("static-type is not set")
		}

		// Get static type from build plan
		var staticType StaticType
		switch rawType {
		case string(HugoType):
			staticType = HugoType
		case string(MdBookType):
			staticType = MdBookType
		default:
			return packit.BuildResult{}, fmt.Errorf("unsupported meta data: %v", plan.Metadata)
		}

		// Get working directory
		workingDir := api.GetWorkingDir(ctx.WorkingDir)

		// Create layers
		staticLayer, err := ctx.Layers.Get("static")
		if err != nil {
			return packit.BuildResult{}, err
		}

		var (
			publicDir = "htdocs"
			args      []string
		)

		switch staticType {
		case HugoType:
			version := os.Getenv(api.HugoVersionEnv)
			if version == "" {
				version = HugoVersion
			}
			if err := installHugo(log, staticLayer, version); err != nil {
				return packit.BuildResult{}, fmt.Errorf("failed to install Hugo: %w", err)
			}
			args = append(args, []string{"--source", ".", "--destination", publicDir, "--minify"}...)
		case MdBookType:
			version := os.Getenv(api.MdBookVersionEnv)
			if version == "" {
				version = MdBookVersion
			}
			if err := installMdBook(log, staticLayer, version); err != nil {
				return packit.BuildResult{}, fmt.Errorf("failed to install mdBook: %w", err)
			}
			args = append(args, []string{"build", ".", "--dest-dir", publicDir}...)
		default:
			return packit.BuildResult{}, fmt.Errorf("unsupported static type: %s", staticType)
		}

		// Set up the build process
		staticLayer.Launch = true
		staticLayer.Build = true
		staticLayer.Cache = true

		os.Setenv("PATH", fmt.Sprintf("%s:%s", filepath.Join(staticLayer.Path, "bin"), os.Getenv("PATH")))

		// Execute the build command during the build phase
		log.Process("Building static site with %s", string(staticType))

		if err := pexec.NewExecutable(string(staticType)).Execute(pexec.Execution{
			Args:   args,
			Dir:    workingDir,
			Stdout: log.ActionWriter,
			Stderr: log.ActionWriter,
		}); err != nil {
			return packit.BuildResult{}, err
		}

		log.Process("Static site built successfully")

		log.Process("Configuring webserver")

		// It's too late to set BP_WEB_SERVER_ROOT (that's a build-time var,
		// but the nginx build is already done here).
		//
		// But the nginx default conf respects "$APP_ROOT" at runtime, so we set that.
		staticLayer.LaunchEnv = packit.Environment{
			"APP_ROOT": publicDir,
		}
		log.EnvironmentVariables(staticLayer)

		return packit.BuildResult{
			Layers: []packit.Layer{
				staticLayer,
			},
		}, nil
	}
}
