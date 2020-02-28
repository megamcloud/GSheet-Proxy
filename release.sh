#!/bin/bash

echo "Prepare Release folder"
rm ./releases -rf
mkdir ./releases

#echo "Build templates assets"
#go-assets-builder templates --package=assets -o infrastructure/assets/assets.go

echo "Build binaries for major platforms"
   GOOS=windows GOARCH=386 go build main.go && mv ./main.exe ./releases/win_scanHub.exe \
&& GOOS=darwin GOARCH=amd64 go build main.go && mv ./main ./releases/darwin_scanHub \
&& GOOS=linux GOARCH=amd64 go build main.go && mv ./main ./releases/linux_scanHub

echo "Copy config"
cp -r ./templates ./releases/
cp -r ./tls ./releases/
cp -r ./public ./releases/
cp config.yaml ./releases/
