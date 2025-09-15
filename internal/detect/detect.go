package detect

import (
	"github.com/hostwithquantum/static-buildpack/api"
	"github.com/hostwithquantum/static-buildpack/internal/meta"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Detect(logs scribe.Emitter) packit.DetectFunc {
	return func(ctx packit.DetectContext) (packit.DetectResult, error) {
		logs.Title("%s %s", ctx.Info.Name, ctx.Info.Version)

		workingDir := api.GetWorkingDir(ctx.CNBPath, ctx.WorkingDir)

		logs.Process("Checking working directory: %s", workingDir)

		finder := NewFinder(logs)
		if err := finder.Find(workingDir); err != nil {
			return packit.DetectResult{}, err
		}

		logs.Detail("Detected static site")
		logs.Detail("Type: %s", finder.GetStaticType())

		webServer := meta.DetectWebServer()

		logs.Process("Select webserver: %s", webServer)

		buildPlanRequirements := []packit.BuildPlanRequirement{
			{
				Name: "static-buildpack",
				Metadata: map[string]any{
					"static-type": string(finder.GetStaticType()),
				},
			},
			{
				Name: webServer,
				Metadata: map[string]any{
					"launch": true,
				},
			},
		}

		if finder.GetStaticType() == "hugo" {
			if meta.NeedsGO(workingDir, logs) {
				buildPlanRequirements = append(buildPlanRequirements, packit.BuildPlanRequirement{
					Name: "go",
					Metadata: map[string]any{
						"build": true,
					},
				})
			}
		}

		if meta.NeedsNPM(workingDir, logs) {
			// require 'node' to trigger node-engine buildpack
			buildPlanRequirements = append(buildPlanRequirements, packit.BuildPlanRequirement{
				Name: "node",
				Metadata: map[string]any{
					"build": true,
				},
			})
			buildPlanRequirements = append(buildPlanRequirements, packit.BuildPlanRequirement{
				Name: "npm",
			})

			// require 'node_modules' and therefor: https://github.com/paketo-buildpacks/npm-install
			buildPlanRequirements = append(buildPlanRequirements, packit.BuildPlanRequirement{
				Name: "node_modules",
				Metadata: map[string]any{
					"build": true,
				},
			})
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{
						Name: "static-buildpack",
					},
				},
				Requires: buildPlanRequirements,
			},
		}, nil
	}
}
