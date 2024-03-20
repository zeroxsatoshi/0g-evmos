#!/bin/bash

pid=`ps -ef | grep evmosd | grep -v grep | awk '{print $2}'`
if [[ "$pid" != "" ]]; then
    echo "Terminate evmosd, pid = $pid"
    kill $pid
fi
