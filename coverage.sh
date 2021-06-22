#!/bin/sh
go test -coverprofile=build/coverage.out . ./cmd/ndocid && go tool cover -html=build/coverage.out
