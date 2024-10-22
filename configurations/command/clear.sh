#!/bin/sh

rm -fr log/*

if [ -r autorun.log ]; then rm autorun.log; fi
if [ -r autorun.err ]; then rm autorun.err; fi
if [ -r PTS ]; then rm PTS; fi
if [ -r syslog ]; then rm syslog; fi
if [ -r syserr ]; then rm syserr; fi
if [ -r stdout ]; then rm stdout; fi

# if [ -r game.core ]; then rm game.core; fi
if [ -r VERSION.txt ]; then rm VERSION.txt; fi
if [ -r DEV_LOG.log ]; then rm DEV_LOG.log; fi
if [ -r mob_count ]; then rm mob_count; fi
if [ -r memory_usage_info.txt ]; then rm memory_usage_info.txt; fi
if [ -r packet_info.txt ]; then rm packet_info.txt; fi
if [ -r p2p_packet_info.txt ]; then rm p2p_packet_info.txt; fi
if [ -r udp_packet_info.txt ]; then rm udp_packet_info.txt; fi
if [ -r ProfileLog ]; then rm ProfileLog; fi