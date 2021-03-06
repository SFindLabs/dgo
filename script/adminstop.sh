#!/bin/bash

PROCESS=`ps -ef|grep dgo_admin|grep -v grep|grep -v PPID|awk '{ print $2}'`
for i in $PROCESS
do
  echo "Kill the dgo_admin process [ $i ]"
  kill -2 $i
done
echo "stop ok"