#!/bin/sh

# compile quest first

# db
./run_db.sh

# Test if we want to force compile quests or if the object folder is empty
if [ "$COMPILE_QUESTS" -eq 1 ] || [ ! -d "{{ .root_dir }}/share/locale/italy/quest/object" ] || [ -z "$(ls -A {{ .root_dir }}/share/locale/italy/quest/object)" ]; then
    # quests (takes some time)
    ./qc.sh
else
    sleep 10
fi

./run.sh