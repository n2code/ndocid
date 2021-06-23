#!/bin/sh
mkdir -p build
GOOS=linux   GOARCH=amd64 go build -trimpath -o ./build/ndocid ./cmd/ndocid
GOOS=windows GOARCH=amd64 go build -trimpath -o ./build/ndocid.exe ./cmd/ndocid
cd build
chmod +x ndocid
sed --in-place '/ndocid -h/q' ../README.md
./ndocid -h >> ../README.md 2>&1
echo '```' >> ../README.md
