package corpus

import "embed"

//go:embed data/*.txt
var Data embed.FS
