#!/bin/bash

# Define the list of target platforms
# PLATFORMS=("darwin-amd64" "linux-386" "linux-amd64" "windows-386" "windows-amd64")
PLATFORMS=("windows" "darwin" "linux" "js")
ARCHITECTURE=("amd64" "amd64" "arm64" "wasm")
# todo find out how to apply these env automatically. NB these are for powershill
# $Env:GOOS = "darwin"; $Env:GOARCH = "amd64"; go build -o mac.dmg main.go 		// mac
# $Env:GOOS = "linux"; $Env:GOARCH = "arm64"; go build -o android.apk main.go 	// android
# $Env:GOOS = "windows"; $Env:GOARCH = "amd64"; go build -o windows.exe main.go // windows
# $Env:GOOS = "js"; $Env:GOARCH = "wasm"; go build -o browser.wasm main.go 		// browser

echo "Starting Building... note not really working properly"

# Get the length of the PLATFORMS array
length=${#PLATFORMS[@]}

# Loop over the indices of the arrays
for (( i=0; i<$length; i++ )); do
    GOOS=${PLATFORMS[$i]}
    GOARCH=${ARCHITECTURE[$i]}
    echo "Building for Platform: $GOOS, Architecture: $GOARCH"

    # Set the GOOS and GOARCH environment variables
    export GOOS
    export GOARCH

    # Determine the binary extension based on the target OS
    BIN_EXT=""
    if [[ "$GOOS" == "windows" ]]; then
    BIN_EXT=".exe"
    elif [[ "$GOOS" == "darwin" ]]; then
    BIN_EXT=""
    elif [[ "$GOOS" == "linux" ]]; then
    BIN_EXT=""
    elif [[ "$GOOS" == "js" ]]; then
    BIN_EXT=".wasm"
    fi

    # Build the executable
    go build -o "$GOOS-$GOARCH$BIN_EXT ../main.go"
done

echo "Build completed."