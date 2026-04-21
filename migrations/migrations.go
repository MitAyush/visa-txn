// Package migrations holds embedded SQL applied when the server starts.
package migrations

import _ "embed"

//go:embed 000001_init.sql
var Init string
