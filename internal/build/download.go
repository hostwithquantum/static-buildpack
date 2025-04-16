package build

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
	HugoBaseURL    = "https://github.com/gohugoio/hugo/releases/download/v%s/hugo_%s_%s-%s.tar.gz"
	MdBookBaseURL  = "https://github.com/rust-lang/mdBook/releases/download/v%s/mdbook-v%s-%s-%s.tar.gz"
	BusyboxBaseURL = "https://busybox.net/downloads/binaries/%s/busybox-%s"
)

func downloadAndInstall(url, destPath string) error {
	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	// Create the destination file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer out.Close()

	// Copy the content
	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Make the file executable
	if err := os.Chmod(destPath, 0755); err != nil {
		return fmt.Errorf("failed to make file executable: %w", err)
	}

	return nil
}

func getHugoURL(version string) string {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}
	return fmt.Sprintf(HugoBaseURL, version, version, runtime.GOOS, arch)
}

func getMdBookURL(version string) string {
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x86_64"
	}
	return fmt.Sprintf(MdBookBaseURL, version, version, runtime.GOOS, arch)
}
