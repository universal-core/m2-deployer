#!/bin/bash

# Function to print usage message
usage() {
    echo "Usage: $0 [-v] [-d directory]"
    exit 1
}

VERBOSE=false
DIRECTORY="."

# Check for command line options
while [ $# -gt 0 ]; do
    case "$1" in
        -v)
            VERBOSE=true
            ;;
        -d)
            if [ -n "$2" ]; then
                DIRECTORY="$2"
                shift
            else
                echo "Error: Directory argument missing for option -d"
                usage
            fi
            ;;
        *)
            echo "Invalid option: $1"
            usage
            ;;
    esac
    shift
done
export VERBOSE

# Find all sanitizer_crash.sh scripts and execute them
find $DIRECTORY -name "sanitizer_crash.sh" -exec bash -c 'dir=$(dirname "$1"); printf "Checking directory: %s\n" "$dir"; if [ "$VERBOSE" = true ]; then bash "$dir/sanitizer_crash.sh" "$dir/autorun.err" -v; else bash "$dir/sanitizer_crash.sh" "$dir/autorun.err"; fi' _ {} \;