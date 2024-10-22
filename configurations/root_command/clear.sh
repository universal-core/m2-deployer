#!/bin/sh
# This script is only here for convience
# It calls all other autorun shellscripts in the
# various subdirectories.

#if [ `id -u` = "0" ]; then
#	echo "This script must not be run as root" 1>&2
#	exit 1
#fi

ROOT=$PWD
{{ range excludeStrings .paths "db" }}
echo "Clearing {{ . }}"
cd $ROOT/{{ . }} && sh clear.sh
{{end}}

sleep 2

