# Match

## Installation

- Install protoc
    - brew install protobuf
- Install Dependencies
    - Install Protocol Buffers (protoc)
    - Install Go plugins for gRPC and Protobuf
```
1. go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
2. go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
    - add `export PATH="$PATH:$(go env GOPATH)/bin"` to ~/.zshrc
