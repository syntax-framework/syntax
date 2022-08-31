package syntax

var defaultLiveReloadPatterns = []string{
	`.*\.(html|htm|js|css|png|jpeg|jpg|gif)$`,
}

type ConfigLiveReload struct {
	Disabled bool
	Interval int
	// milliseconds to wait before sending live reload events to the browser. Defaults to `0`.
	Debounce int
	// A list of patterns to trigger the live reloading. This option is required to enable any live reloading
	Patterns []string
	// the endpoint of the live reload SSE event. Defaults to `dev.livereload`.
	Endpoint string
	// If true, CSS changes will trigger a full page reload like other asset types instead of
	// the default hot reload. Useful when class names are determined at runtime, for example
	// when working with CSS modules. Defaults to false.
	ReloadPageOnCss bool
}

type Config struct {
	Dev        bool
	LiveReload ConfigLiveReload
}
