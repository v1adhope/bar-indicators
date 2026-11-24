package queires

import _ "embed"

//go:embed schema.sql
var Schema string

//go:embed record.sql
var Record string

//go:embed get_latest_success.sql
var GetLatestSuccess string
