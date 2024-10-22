#!/bin/sh
# Installs the needed parts
# (Executes child scripts)

#if [ `id -u` = "0" ]; then
#	echo "This script must not be run as root" 1>&2
#	exit 1
#fi

ROOT=$PWD

chmod u=rwx,g=r,o=r $ROOT/share/bin/game
chmod u=rwx,g=r,o=r $ROOT/share/bin/db
# chmod u=rwx,g=r,o=r $ROOT/share/bin/vrunner

find . -type d -exec chmod u=rwx,g=rx,o=rx {} \;
find . -name "*.sh" -exec chmod u=rwx,g=r,o=r {} \;

{{range .paths}}echo "Installing {{ . }}"
cd $ROOT/{{ . }} && sh install.sh
{{end}}

exit 0