#!/bin/bash

# Path to your supervisord configuration file
CONFIG_FILE="./vrunner.config"

# Function to stop all running supervisord processes
stop_supervisord() {
  echo "Stopping all running supervisord processes..."
  # Get the list of supervisord PIDs, excluding this script's PID and any grep commands
  PIDS=$(ps aux | grep '[s]upervisord' | awk '{print $2}' | grep -v $$)
  if [ -n "$PIDS" ]; then
    # Kill all running supervisord processes
    kill -9 $PIDS
    echo "All supervisord processes have been stopped."
  else
    echo "No running supervisord processes found."
  fi
}

# Check if any supervisord process is running
if ps aux | grep '[s]upervisord -c' > /dev/null; then
  # Stop all running supervisord processes
  stop_supervisord
fi

# Start supervisord with the specified configuration file
echo "Starting supervisord with the configuration file: $CONFIG_FILE"
supervisord -c "$CONFIG_FILE"

# Confirm that supervisord has started
if ps aux | grep '[s]upervisord' > /dev/null; then
  echo "Supervisord started successfully."
else
  echo "Failed to start supervisord."
fi
