#!/bin/bash
#run tests in service package
go test ./internal/service 
go env -w GOBIN="/Users/$USER/Software/release/Go/bin"
go env -w GOMODCACHE=/Users/$USER/Software/Go/pkg/mod
go env -w GOPATH="/Users/$User/Software/Go"
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc --go_out=internal/protos --go_opt=paths=source_relative \
       --go-grpc_out=internal/protos --go-grpc_opt=paths=source_relative \
       resources/protos/*.proto


protoc --go_out=internal \
       --go-grpc_out=internal \
       resources/protos/*.proto       

docker build -t processor . --no-cache --progress=plain       