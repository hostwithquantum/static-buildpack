package build

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/paketo-buildpacks/packit/v2/scribe"
	"github.com/paketo-buildpacks/packit/v2/vacation"
)

const (
	// https://github.com/gohugoio/hugo/releases/download/v0.146.5/hugo_0.146.5_linux-amd64.tar.gz
	HugoBaseURL = "https://github.com/gohugoio/hugo/releases/download/v{{.Version}}/hugo_{{.Version}}_{{.OS}}-{{.Arch}}.tar.gz"

	// https://github.com/rust-lang/mdBook/releases/download/v0.4.48/mdbook-v0.4.48-x86_64-unknown-linux-gnu.tar.gz
	MdBookBaseURL = "https://github.com/rust-lang/mdBook/releases/download/v{{.Version}}/mdbook-v{{.Version}}-{{.Arch}}-{{.OS}}.tar.gz"
)

func downloadAndInstall(log scribe.Emitter, url, destPath string) error {
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed request file file (%s): %w", url, err)
	}
	defer resp.Body.Close()

	log.Detail("Downloaded: %s", url)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file (%s): status code %d", url, resp.StatusCode)
	}

	zip := vacation.NewGzipArchive(resp.Body).StripComponents(0)
	if err := zip.Decompress(destPath); err != nil {
		return err
	}

	log.Detail("Uncompressed to: %s", destPath)

	return nil
}

func getHugoURL(version string) string {
	return tprintf(HugoBaseURL, map[string]any{
		"Version": version,
		"OS":      "linux",
		"Arch":    "amd64",
	})
}

func getMdBookURL(version string) string {
	return tprintf(MdBookBaseURL, map[string]any{
		"Version": version,
		"Arch":    "x86_64",
		"OS":      "unknown-linux-gnu",
	})
}

// when you don't like fmt.Spintf()
func tprintf(tmpl string, data map[string]any) string {
	t := template.Must(template.New("url").Parse(tmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return ""
	}
	return buf.String()
}
