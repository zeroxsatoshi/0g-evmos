#!/bin/bash

pid=`ps -ef | grep evmos | grep -v grep | awk '{print $2}'`
if [[ "$pid" != "" ]]; then
    echo "Terminate evmos"
    kill $pid
fi
