// Package public embeds public web assets.
package public

import (
	"embed"
)

//go:embed css/* js/* login/* all:signup
var HtmlEFS embed.FS
