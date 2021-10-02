package client

const (
	FlibustaHostEnvKey = "FLIBUSTA_HOST"
	Fb2                = "fb2"
	Epub               = "epub"
	Mobi               = "mobi"
)

var (
	FlibustaMirrors = []string{
		"flibusta.is",
		"flibusta.site",
		"flibustahezeous3.onion",
	}
	validFormats = []string{Fb2, Epub, Mobi}
)
