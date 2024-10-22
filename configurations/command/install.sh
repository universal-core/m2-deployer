#!/bin/sh
# This simple script installs the needed parts
# and files to run the server properly.
# Do not change anything here unless you know
# what you are doing.

{{/* ------------------------------------ DB INSTALL */}}
{{ if eq .hostname "db" }}
# Delete db if already present
rm -f {{ .hostname }} locale item_proto.txt item_proto_test.txt mob_proto.txt mob_proto_test.txt item_names.txt mob_names.txt

# Create the symbolic links
ln -s {{ .root_dir }}/share/conf/item_proto.txt item_proto.txt
#ln -s {{ .root_dir }}/share/conf/item_proto_test.txt item_proto_test.txt
ln -s {{ .root_dir }}/share/conf/mob_proto.txt mob_proto.txt
#ln -s {{ .root_dir }}/share/conf/mob_proto_test.txt mob_proto_test.txt
ln -s {{ .root_dir }}/share/conf/item_names.txt item_names.txt
ln -s {{ .root_dir }}/share/conf/mob_names.txt mob_names.txt
ln -s {{ .root_dir }}/share/bin/db {{ .hostname }}

ln -s {{ .root_dir }}/share/locale locale

if [ ! -d log ]; then mkdir log; fi
if [ ! -d cores ]; then mkdir cores; fi

# chmod u=rwx,g=r,o=r vrunner
chmod u=rwx,g=r,o=r {{ .hostname }}
chmod u=rw,g=r,o= *.txt
chmod u=rwx,g=rx,o=rx log
chmod u=rwx,g=rx,o=rx cores
chmod u=r,g=r,o=r log/*/syslog.* >/dev/null 2>&1
{{/* ------------------------------------ AUTH INSTALL */}}
{{ else if eq .hostname "auth" }}
# Delete the files if they're already present
rm -f {{ .hostname }} data locale common.xml

# Create the symbolic links
ln -s {{ .root_dir }}/share/bin/game {{ .hostname }}
ln -s {{ .root_dir }}/share/data data
ln -s {{ .root_dir }}/share/locale locale
ln -s {{ .root_dir }}/share/conf/common.xml common.xml

if [ ! -d log ]; then mkdir log; fi
if [ ! -d cores ]; then mkdir cores; fi

# crypto
if [ ! -d crypto ]; then mkdir crypto; fi
if [ ! -d crypto/index ]; then mkdir crypto/index; fi
if [ ! -d crypto/key ]; then mkdir crypto/index; fi
touch crypto/crypto.lst
echo "0" >> crypto/crypto.ver

chmod u=rwx,g=r,o=r {{ .hostname }}
chmod u=rw,g=r,o= server.xml
chmod u=rwx,g=rx,o=rx log
chmod u=rwx,g=rx,o=rx cores
chmod u=rwx,g=rx,o=rx crypto
chmod u=rwx,g=rx,o=rx crypto/index
chmod u=rwx,g=rx,o=rx crypto/key
chmod u=rw,g=r,o= crypto/crypto.lst
chmod u=rw,g=r,o= crypto/crypto.ver
chmod u=r,g=r,o=r log/*/syslog.* >/dev/null 2>&1
{{/* ------------------------------------ GAME INSTALL */}}
{{ else }}
# Delete the files if they're already present
rm -f {{ .hostname }} data locale cmd.xml common.xml 3N.mhe HShield.dat HSPub.key metin2client.hsb

# Create the symbolic links
ln -s {{ .root_dir }}/share/bin/game {{ .hostname }}
ln -s {{ .root_dir }}/share/data data
ln -s {{ .root_dir }}/share/locale locale
ln -s {{ .root_dir }}/share/conf/cmd.xml cmd.xml
ln -s {{ .root_dir }}/share/conf/common.xml common.xml

if [ ! -d mark ]; then mkdir mark; fi
if [ ! -d log ]; then mkdir log; fi
if [ ! -d cores ]; then mkdir cores; fi

chmod u=rwx,g=r,o=r {{ .hostname }}
chmod u=rw,g=r,o= server.xml
chmod u=rwx,g=rx,o=rx log
chmod u=rwx,g=rx,o=rx cores
chmod u=r,g=r,o=r log/*/syslog.* >/dev/null 2>&1
{{ end }}