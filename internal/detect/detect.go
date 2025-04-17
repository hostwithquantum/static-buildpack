package detect

import (
	"github.com/hostwithquantum/static-buildpack/internal/meta"
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

		staticType := finder.GetStaticType()

		webServer := meta.DetectWebServer()
		publicDir := meta.DetectHtDocs()

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{
						Name: string(staticType),
					},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: ctx.Info.Name,
						Metadata: map[string]any{
							"type":  string(staticType),
							"build": true,
						},
					},
					{
						Name: "paketo-buildpacks/" + webServer,
						Metadata: map[string]any{
							"launch": true,
							"env": map[string]string{
								"BP_WEB_SERVER_ROOT": publicDir,
							},
						},
					},
				},
			},
		}, nil
	}
}
