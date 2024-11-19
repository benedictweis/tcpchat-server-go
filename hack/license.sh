#!/usr/bin/env bash

# Check if the correct number of arguments are provided
if [ "$#" -lt 2 ] || [ "$#" -gt 3 ]; then
    echo "Usage: $0 <license_file> <file_pattern> [fix]"
    exit 1
fi

LICENSE_FILE="$1"
FILE_PATTERN="$2"
FIX_MODE="$3"

# Check if the license file exists
if [ ! -f "$LICENSE_FILE" ]; then
    echo "Error: License file '$LICENSE_FILE' not found."
    exit 1
fi

# Read the license file into a variable
LICENSE_CONTENT=$(<"$LICENSE_FILE")

# Get the number of lines in the license file
LICENSE_LINES=$(wc -l < "$LICENSE_FILE")

# Flag to track if any file is missing the license
MISSING_LICENSE=0

# Iterate over the files matching the pattern
for FILE in $(find . -type f -name "$FILE_PATTERN"); do
    # Extract the top lines from the file to compare with the license
    FILE_TOP_CONTENT=$(head -n "$LICENSE_LINES" "$FILE")

    # Check if the file top content matches the license exactly
    if [ "$FILE_TOP_CONTENT" != "$LICENSE_CONTENT" ]; then
        if [ "$FIX_MODE" == "true" ]; then
            # Prepend the license content to the file with a newline after the license
            echo "Adding license to '$FILE'."
            { echo "$LICENSE_CONTENT"; echo ""; cat "$FILE"; } > "${FILE}.tmp" && mv "${FILE}.tmp" "$FILE"
        else
            echo "License missing or incorrect in '$FILE'."
            MISSING_LICENSE=1
        fi
    fi
done

# Exit with 1 if any file is missing the license and fix mode is not enabled
if [ "$MISSING_LICENSE" -eq 1 ] && [ "$FIX_MODE" != "true" ]; then
    exit 1
else
    exit 0
fi
