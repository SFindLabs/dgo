#!/bin/bash
host="127.0.0.1"
db_user="xxxxx_user"
db_password="G5xxxxxxxxxdd"
db_name="xxxxx"

#Backup Dir
backup_dir="/data/backupsql/xxxxxx"
mkdir -p $backup_dir

#Time Format
time=$(date +"%Y%m%d%H%M%S")

#Execute Backup
mysqldump -h$host -u$db_user -p$db_password $db_name --default-character-set=utf8mb4 | gzip > "$backup_dir/$db_name"-"$time.sql.gz"
echo "down ok"
ls -l -h  /data/backupsql/xxxxxxxx/