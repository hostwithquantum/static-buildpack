package main

import (
	"os"

	"github.com/hostwithquantum/static-buildpack/internal/detect"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

func main() {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))
	packit.Detect(detect.Detect(logEmitter))
}
