#!/bin/sh
# Shuts each process down
# (Executes child scripts)

#if [ `id -u` = "0" ]; then
#	echo "This script must not be run as root" 1>&2
#	exit 1
#fi

ROOT=$PWD

echo "Stopping {{ index (filterStrings .paths "db") 0 }}"
cd $ROOT/{{ index (filterStrings .paths "db") 0 }} && sh shutdown.sh

# Check processes are not still running
VRUNNER_RUNNING=$(ps -c | grep -c vrunner)

if [ $VRUNNER_RUNNING != '0' ]; then
	echo ""
	echo "Still running proccess found:"
	
	PROC=$(ps xco pid,command | grep vrunner)
	echo -e "\e[1;31m$PROC\e[0m"
fi
