#!/bin/bash

pid=`ps -ef | grep evmos | grep -v grep | awk '{print $2}'`
if [[ "$pid" = "" ]]; then
    nohup evmosd start --home evmosd-prod &
else
    echo "evmos already started"
fi
