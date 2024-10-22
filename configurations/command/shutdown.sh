#!/bin/sh
supervisorctl -c {{ .root_dir }}/vrunner.config stop {{ .hostname }}