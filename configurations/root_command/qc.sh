#!/bin/bash

# Function to print the help message
print_help() {
    echo "Usage: ./qc.sh [--f <filepath>]"
    exit 1
}

# Check if --help is present in the arguments
if [ " $* " == *" --help "* ]; then
    print_help
fi

# Initialize variables
path=""

# Parse arguments
while [[ $# -gt 0 ]]; do
    key="$1"

    case $key in
        --f)
        path="$2"
        shift
        shift
        ;;
        *)
        print_help
        ;;
    esac
done

chmod +x share/bin/qc
cd share

qc_args="--cwd=locale/italy/quest --out_forge_game=./forge --out_forge_db=./forge --out_forge_auth=./forge --out_quest=./object"

if [ -z "$path" ]; then
    echo "Compiling all quests"
    ./bin/qc $qc_args --clean --do_file=qc/quest_list
else
    echo "Compiling quest: $path"
    ./bin/qc $qc_args --files="$path"
fi