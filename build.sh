#!/bin/bash

# Linux host machine default
output="deployer"
# Check the value of the OSTYPE environment variable
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
  # Windows host machine
  output="deployer.exe"
fi

cd ./src/
go mod tidy
go build -o ../"$output" .
