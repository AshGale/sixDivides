#!/bin/bash

# Define the lists
PLATFORMS=("windows" "js")
ARCHITECTURES=("amd64" "wasm")

echo "Starting Building..."

# Get the length of the arrays (assuming both arrays are of the same length)
length=${#PLATFORMS[@]}

# Loop over the list of platforms, and build
for ((i=0; i<length; i++)); do
    GOOS=${PLATFORMS[i]}
    GOARCH=${ARCHITECTURES[i]}

    # Determine the binary extension based on the platform
    BIN_EXT=""
    case "$GOOS" in
        "windows")
            BIN_EXT=".exe"
            ;;
        "js")
            BIN_EXT=".wasm"
            ;;
    esac

    echo "  Platform: $GOOS, Architecture: $GOARCH, Extension: $BIN_EXT"

    go build -o "builds/$GOOS-$GOARCH$BIN_EXT" main.go
done

echo "Finished Building!"

# Define the source file and destination directory
SOURCE_FILE="./builds/js-wasm.wasm"
DEST_DIRECTORY="../server/static/wasm"
DEST_FILE="$DEST_DIRECTORY/sixDivides.wasm"

# Ensure the destination directory exists
if [ ! -d "$DEST_DIRECTORY" ]; then
    echo "Destination directory does not exist. Creating directory..."
    mkdir -p "$DEST_DIRECTORY"
fi

# Check if the source file exists
if [ ! -f "$SOURCE_FILE" ]; then
    echo "Source file $SOURCE_FILE does not exist."
    exit 1
fi

# Copy the file with the new name to the destination directory
cp "$SOURCE_FILE" "$DEST_FILE"

# Check if the copy was successful
if [ -f "$DEST_FILE" ]; then
    echo "File copied successfully to $DEST_FILE."
else
    echo "Failed to copy file."
    exit 1
fi

