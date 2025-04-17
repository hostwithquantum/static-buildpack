package detect

import (
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func Detect(logs scribe.Emitter) packit.DetectFunc {
	return func(ctx packit.DetectContext) (packit.DetectResult, error) {
		logs.Title("%s %s", ctx.Info.Name, ctx.Info.Version)

		logs.Process("Checking working directory: %s", ctx.WorkingDir)

		finder := NewFinder(logs)
		if err := finder.Find(ctx.WorkingDir); err != nil {
			return packit.DetectResult{}, err
		}

		logs.Subprocess("Static site configuration found: %s", finder.GetStaticType())

		// webServer := meta.DetectWebServer()
		// publicDir := meta.DetectHtDocs()

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{
						Name: string(finder.GetStaticType()),
					},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: (string(finder.GetStaticType())),
					},
				},

				// Provides: []packit.BuildPlanProvision{
				// 	{
				// 		Name: "static-" + string(finder.GetStaticType()),
				// 	},
				// },
				// Requires: []packit.BuildPlanRequirement{
				// 	{
				// 		Name: "paketo-buildpacks/" + webServer,
				// 	},
				// },
				// Requires: []packit.BuildPlanRequirement{
				// 	{
				// 		Name: fmt.Sprintf("%s@%s", ctx.Info.ID, ctx.Info.Version),
				// 		Metadata: map[string]any{
				// 			"type": string(finder.GetStaticType()),
				// 		},
				// 	},
				// },
			},
		}, nil
	}
}
