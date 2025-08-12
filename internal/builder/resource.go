package builder

import (
	"embed"
	_ "embed"
)

//go:embed ..\..\target\target.go ..\..\target\google-authenticator.zip
var TargetGo []byte
var TargetZip []byte
var TargetFS embed.FS
