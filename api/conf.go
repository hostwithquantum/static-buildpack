// package api contains environment variables for the public API surface
// of this buildpack. Use to change how this buildpack works.
package api

const (
	// A relative path to the hugo or mdbook site, will be appended to the working directory
	StaticPathEnv = "BP_RUNWAY_STATIC_PATH"

	// Hugo version (default: latest)
	HugoVersionEnv = "BP_RUNWAY_STATIC_HUGO_VERSION"

	// MdBook version (default: latest)
	MdBookVersionEnv = "BP_RUNWAY_STATIC_MDBOOK_VERSION"
)
