#!/bin/bash
#run tests in service package
go test ./internal/service 
go env -w GOBIN="/Users/$USER/Software/release/Go/bin"
go env -w GOMODCACHE=/Users/$USER/Software/Go/pkg/mod
go env -w GOPATH="/Users/$User/Software/Go"