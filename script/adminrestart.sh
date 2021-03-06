#!/bin/bash
./script/adminstop.sh
sleep 1s

nohup ./dgo_admin -t admin >> log/admin.log 2>&1 &
echo "restart ok"

sleep 1s
ps -ef|grep dgo_admin|grep -v grep