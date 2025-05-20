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

		workingDir := api.GetWorkingDir(ctx.WorkingDir)

		logs.Process("Checking working directory: %s", workingDir)

		finder := NewFinder(logs)
		if err := finder.Find(workingDir); err != nil {
			return packit.DetectResult{}, err
		}

		logs.Detail("Detected static site")
		logs.Detail("Type: %s", finder.GetStaticType())

		webServer := meta.DetectWebServer()

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{
						Name: "static-buildpack",
					},
				},
				Requires: []packit.BuildPlanRequirement{
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
				},
			},
		}, nil
	}
}
