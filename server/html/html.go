package html

import (
	"embed"
)

//go:embed view
var View embed.FS

//go:embed static
var Static embed.FS
