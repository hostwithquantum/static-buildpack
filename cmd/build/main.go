package main

import (
	"os"

	"github.com/hostwithquantum/static-buildpack/internal/build"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func main() {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))
	packit.Build(build.Build(logEmitter))
}
