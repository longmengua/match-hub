package enum

// go install golang.org/x/tools/cmd/stringer@v0.13.0
// go generate

//go:generate stringer -type=Code
type Code int

const (
	CodeOk           = Code(102000)
	CodeExists       = Code(504001)
	CodeNotFound     = Code(504040)
	CodeParamInvalid = Code(504000)
)
