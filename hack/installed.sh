#!/usr/bin/env bash

# Check if a software is installed
check_software_installed() {
    # Argument: Software name
    local software=$1

    # Check if the software is available in the PATH
    if ! command -v "$software" &> /dev/null; then
        echo "Error: '$software' is not installed, please install it!" >&2
        return 1
    fi

    return 0
}

# Usage
if [ $# -eq 0 ]; then
    echo "Usage: $0 <software_name>" >&2
    exit 2
fi

# Check for the software
check_software_installed "$1"
exit $?
