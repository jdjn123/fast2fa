package builder

import (
	_ "embed"
)

//go:embed ..\..\target\target.go
var TargetGo []byte
