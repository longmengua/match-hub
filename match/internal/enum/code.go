package enum

// go install golang.org/x/tools/cmd/stringer@v0.13.0
// go generate

//go:generate stringer -type=Code
type Code int

const (
	Ok            = Code(200001)
	InternalError = Code(500001)
	InvalidParams = Code(400001)
)
