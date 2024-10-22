#!/bin/bash

# Prefix stdout and stderr
exec 1> >( perl -ne 'use Time::HiRes qw(time); use POSIX qw( strftime ); $time=time; $microsecs = ($time - int($time)) * 1e3; $| = 1; printf( "'"${SUPERVISOR_PROCESS_NAME}"' | %s,%03.0f %s", strftime("%Y-%m-%d %H:%M:%S", gmtime($time)), $microsecs, $_);' >&1)
exec 2> >( perl -ne 'use Time::HiRes qw(time); use POSIX qw( strftime ); $time=time; $microsecs = ($time - int($time)) * 1e3; $| = 1; printf( "'"${SUPERVISOR_PROCESS_NAME}"' | %s,%03.0f %s", strftime("%Y-%m-%d %H:%M:%S", gmtime($time)), $microsecs, $_);' >&2)

exec "$@"