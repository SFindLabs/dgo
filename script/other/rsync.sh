#!/bin/bash

time=$(date +"%Y%m%d%H%M%S")
backup_dir="/data/backup/xxxxxxx/$time/"
mkdir -p $backup_dir
cp -rf /data/golang/xxxxxxx/{xxxxxxx,conf,script} $backup_dir
rsync -ravzp -e 'ssh -p 22' root@192.0.0.0:/data/golang/gopath/src/xxxxxxx/xxxxxxx /data/golang/xxxxxxx/
