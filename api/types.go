// package api contains environment variables for the public API surface
// of this buildpack. Use to change how this buildpack works.
package api

const (
	// A relative path to the hugo or mdbook site, will be appended to the working directory
	StaticPathEnv = "BP_RUNWAY_STATIC_PATH"

	// Hugo version
	HugoVersionEnv = "BP_RUNWAY_STATIC_HUGO_VERSION"

	// MdBook version
	MdBookVersionEnv = "BP_RUNWAY_STATIC_MDBOOK_VERSION"
)

type BuildpackTOML struct {
	Buildpack struct {
		ID       string `toml:"id"`
		Name     string `toml:"name"`
		Version  string `toml:"version"`
		Homepage string `toml:"homepage"`
	} `toml:"buildpack"`
	Metadata struct {
		Configurations []struct {
			Name    string `toml:"name"`
			Default string `toml:"default"`
		} `toml:"configurations"`
	} `toml:"metadata"`
}
