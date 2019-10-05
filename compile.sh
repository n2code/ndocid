#!/bin/sh
GOOS=linux   GOARCH=amd64 go build -o ndocid
GOOS=windows GOARCH=amd64 go build -o ndocid.exe
