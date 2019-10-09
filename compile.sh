#!/bin/sh
GOOS=linux   GOARCH=amd64 go build -o ndocid
GOOS=windows GOARCH=amd64 go build -o ndocid.exe
chmod +x ndocid
sed --in-place '/ndocid -h/q' README.md
./ndocid -h >> README.md 2>&1
echo '```' >> README.md
