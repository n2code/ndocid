#!/bin/sh
go test -coverprofile=cov.out && go tool cover -html=cov.out && rm cov.out
