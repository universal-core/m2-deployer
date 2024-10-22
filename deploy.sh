#!/bin/bash

# Function to display help message
function display_help {
    echo "Usage: $0 --out [out_path] --config [configurations] --channels [number_of_channels]"
    exit 1
}

# Initialize variables
out_path=""
configurations=""
channels=""

# Parse command-line arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --out) out_path="$2"; shift ;;
        --config) configurations="$2"; shift ;;
        --channels) channels="$2"; shift ;;
        *) echo "Unknown parameter passed: $1"; display_help ;;
    esac
    shift
done

# Check if all required arguments are provided
if [[ -z "$out_path" || -z "$configurations" || -z "$channels" ]]; then
    display_help
fi


# Execute the deployer command with provided arguments
./deployer --out "$out_path" --config "deploy_base.yaml,$configurations" --channels "$channels"

# Set bash permissions
cd "$out_path"

find . -type f -name "*.sh" -exec chmod +x {} \;
