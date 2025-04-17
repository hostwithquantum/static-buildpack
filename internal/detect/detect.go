package detect

import (
	"os"

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

		// Determine which web server to use
		webServer := os.Getenv("BP_WEB_SERVER")
		if webServer == "" {
			webServer = "nginx" // Default to nginx
		}

		// determine the public dir
		var publicDir = "htdocs"
		if os.Getenv("BP_WEB_SERVER_ROOT") != "" {
			publicDir = os.Getenv("BP_WEB_SERVER_ROOT")
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{
						Name: string(staticType),
					},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: string(staticType),
						Metadata: map[string]any{
							"build": true,
						},
					},
					{
						Name: "paketo-buildpacks/web-servers",
						Metadata: map[string]any{
							"launch": true,
							"env": map[string]string{
								"BP_WEB_SERVER":      "nginx",
								"BP_WEB_SERVER_ROOT": publicDir,
							},
						},
					},
				},
			},
		}, nil
	}
}
