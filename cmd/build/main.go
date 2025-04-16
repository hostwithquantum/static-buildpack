package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hostwithquantum/static-buildpack/internal/build"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

type (
	StaticType string
	WebServer  string
)

const (
	HugoType   StaticType = "hugo"
	MdBookType StaticType = "mdbook"

	// Web server types
	HttpdServer WebServer = "httpd"
	NginxServer WebServer = "nginx"

	// Environment variables
	StaticPathEnv    = "BP_RUNWAY_STATIC_PATH"
	HugoVersionEnv   = "BP_RUNWAY_STATIC_HUGO_VERSION"
	MdBookVersionEnv = "BP_RUNWAY_STATIC_MDBOOK_VERSION"
	WebServerEnv     = "BP_RUNWAY_STATIC_WEB_SERVER"
)

func run(logEmitter scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logEmitter.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		// Get static type from build plan
		staticType, ok := context.Plan.Entries[0].Metadata["type"].(StaticType)
		if !ok {
			return packit.BuildResult{}, fmt.Errorf("invalid static type in build plan")
		}

		// Get working directory
		workingDir := context.WorkingDir
		if path := os.Getenv(StaticPathEnv); path != "" {
			workingDir = filepath.Join(workingDir, path)
		}

		// Create layers
		staticLayer, err := context.Layers.Get("static")
		if err != nil {
			return packit.BuildResult{}, err
		}

		// Download and install the appropriate tool
		var buildCmd string
		var toolName string
		switch staticType {
		case HugoType:
			version := os.Getenv(HugoVersionEnv)
			if version == "" {
				version = "latest"
			}
			if err := build.InstallHugo(staticLayer, version); err != nil {
				return packit.BuildResult{}, fmt.Errorf("failed to install Hugo: %w", err)
			}
			toolName = "hugo"
			buildCmd = build.GetStaticToolPath(staticLayer, toolName)
		case MdBookType:
			version := os.Getenv(MdBookVersionEnv)
			if version == "" {
				version = "latest"
			}
			if err := build.InstallMdBook(staticLayer, version); err != nil {
				return packit.BuildResult{}, fmt.Errorf("failed to install mdBook: %w", err)
			}
			toolName = "mdbook"
			buildCmd = build.GetStaticToolPath(staticLayer, toolName)
		default:
			return packit.BuildResult{}, fmt.Errorf("unsupported static type: %s", staticType)
		}

		// Set up the build process
		staticLayer.Launch = true
		staticLayer.Build = true
		staticLayer.Cache = true

		// Execute the build command during the build phase
		logEmitter.Process("Building static site with %s", toolName)

		cmd := exec.Command(buildCmd, "build")
		cmd.Dir = workingDir
		cmd.Env = append(os.Environ(), "PATH="+filepath.Join(staticLayer.Path, "bin")+":"+os.Getenv("PATH"))

		output, err := cmd.CombinedOutput()
		if err != nil {
			logEmitter.Detail("Build output: %s", string(output))
			return packit.BuildResult{}, fmt.Errorf("failed to build static site: %w", err)
		}

		logEmitter.Detail("Build output: %s", string(output))
		logEmitter.Process("Static site built successfully")

		// Determine which web server to use
		webServer := WebServer(os.Getenv(WebServerEnv))
		if webServer == "" {
			webServer = NginxServer // Default to nginx
		}

		// Set up the web server process
		var process packit.Process
		outputDir := getOutputDir(staticType, workingDir)

		switch webServer {
		case HttpdServer:
			process = packit.Process{
				Type:    "web",
				Command: "httpd",
				Args:    []string{"-f", "-v", "-p", "8080", "-h", outputDir},
				Direct:  true,
			}
		case NginxServer:
			process = packit.Process{
				Type:    "web",
				Command: "nginx",
				Args:    []string{"-g", "daemon off;"},
				Direct:  true,
			}
		default:
			return packit.BuildResult{}, fmt.Errorf("unsupported web server: %s", webServer)
		}

		// Add the web server as a requirement
		requirements := []packit.BuildPlanRequirement{
			{
				Name: string(webServer),
				Metadata: map[string]interface{}{
					"version": "0.7.0",
				},
			},
		}

		entries := make([]packit.BuildpackPlanEntry, len(requirements))
		for i, req := range requirements {
			metadata, ok := req.Metadata.(map[string]interface{})
			if !ok {
				metadata = make(map[string]interface{})
			}
			entries[i] = packit.BuildpackPlanEntry{
				Name:     req.Name,
				Metadata: metadata,
			}
		}

		return packit.BuildResult{
			Layers: []packit.Layer{
				staticLayer,
			},
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{
					process,
				},
			},
			Plan: packit.BuildpackPlan{
				Entries: entries,
			},
		}, nil
	}
}

// getOutputDir returns the directory where the static site is built
func getOutputDir(staticType StaticType, workingDir string) string {
	switch staticType {
	case HugoType:
		return filepath.Join(workingDir, "public")
	case MdBookType:
		return filepath.Join(workingDir, "book")
	default:
		return workingDir
	}
}

func main() {
	logEmitter := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))
	packit.Build(run(logEmitter))
}
