#!/bin/bash

pid=`ps -ef | grep evmosd | grep -v grep | awk '{print $2}'`
if [[ "$pid" = "" ]]; then
    nohup evmosd start --home ~/evmosd-prod > ~/stdout.evmosd.log 2>~/stderr.evmosd.log &
else
    echo "evmosd already started, pid = $pid"
fi
