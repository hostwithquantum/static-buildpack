package meta

import (
	"os"
	"slices"
)

// Determine which web server to use
func DetectWebServer() string {
	webServer := os.Getenv("BP_WEB_SERVER")
	if slices.Contains([]string{"httpd", "nginx"}, webServer) {
		return webServer
	}
	return "nginx"
}
