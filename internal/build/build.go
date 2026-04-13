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
)

// resolveVersion determines which version to use, preferring:
// 1) explicit env var set by the user
// 2) cached version from a previous build
// 3) Default from buildpack.toml
func ResolveVersion(log scribe.Emitter, layer packit.Layer, cnbPath, envVar string) string {
	// If the user explicitly set the env var, always use that
	if version := os.Getenv(envVar); version != "" {
		return version
	}

	// Check if we have a cached version from a previous build
	if cached, ok := layer.Metadata["version"].(string); ok && cached != "" {
		log.Process("Reusing previously installed version %s (set %s to override)", cached, envVar)
		return cached
	}

	// Fall back to buildpack.toml default
	return api.GetDefault(cnbPath, envVar)
}

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
		workingDir := api.GetWorkingDir(ctx.CNBPath, ctx.WorkingDir)

		// Create layers
		staticLayer, err := ctx.Layers.Get("static")
		if err != nil {
			return packit.BuildResult{}, err
		}

		var (
			publicDir = "htdocs"
			args      []string
			version   string
		)

		switch staticType {
		case HugoType:
			version = ResolveVersion(log, staticLayer, ctx.CNBPath, api.HugoVersionEnv)
			if version == "" {
				return packit.BuildResult{}, fmt.Errorf("no Hugo version specified and no default found in buildpack.toml")
			}
			args = append(args, []string{"--source", ".", "--destination", publicDir, "--minify"}...)
		case MdBookType:
			version = ResolveVersion(log, staticLayer, ctx.CNBPath, api.MdBookVersionEnv)
			if version == "" {
				return packit.BuildResult{}, fmt.Errorf("no MdBook version specified and no default found in buildpack.toml")
			}
			args = append(args, []string{"build", ".", "--dest-dir", publicDir}...)
		default:
			return packit.BuildResult{}, fmt.Errorf("unsupported static type: %s", staticType)
		}

		// Check if cached layer already has the right version installed
		cachedVersion, _ := staticLayer.Metadata["version"].(string)
		if cachedVersion == version {
			log.Process("Reusing cached %s %s", string(staticType), version)
		} else {
			// Version changed or first build — reset and reinstall
			staticLayer, err = staticLayer.Reset()
			if err != nil {
				return packit.BuildResult{}, fmt.Errorf("failed to reset layer: %w", err)
			}

			switch staticType {
			case HugoType:
				if err := installHugo(log, staticLayer, version); err != nil {
					return packit.BuildResult{}, fmt.Errorf("failed to install Hugo: %w", err)
				}
			case MdBookType:
				if err := installMdBook(log, staticLayer, version); err != nil {
					return packit.BuildResult{}, fmt.Errorf("failed to install mdBook: %w", err)
				}
			}

			staticLayer.Metadata = map[string]any{
				"version": version,
			}
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
