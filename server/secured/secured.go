// Package secured embeds assets that require authentication for access.

package secured

import (
	"embed"
)

//go:embed all:dist
var SecuredEFS embed.FS
