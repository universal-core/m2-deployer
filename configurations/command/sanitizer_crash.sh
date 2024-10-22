#!/bin/sh

# Function to print usage message
usage() {
    echo "Usage: $0 <filename> [-v]"
    exit 1
}

# Check if filename is provided
if [ -z "$1" ]; then
    usage
fi

VERBOSE=false
filename="$1"
shift

# Check for command line options
while [ $# -gt 0 ]; do
    case "$1" in
        -v)
            VERBOSE=true
            ;;
        *)
            echo "Invalid option: $1"
            usage
            ;;
    esac
    shift
done

# Check if file exists
if [ ! -f "$filename" ]; then
    echo "File $1 not found"
    usage
fi

# Count occurrences of "ABORTING"
aborting_count=$(grep -c "ABORTING" "$filename")
echo "Crashes occurred: $aborting_count"

# Print verbose output
if [ "$VERBOSE" = true ]; then
    # Use awk to print chunks between ==ERROR and ==ABORTING
    awk '/==ERROR/ && !/LeakSanitizer/{p=1} p; /==ABORTING/{print "\n-----------------------\n"; p=0}' "$filename"
fi
