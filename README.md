# Mysql Backup Cron

## Install

`go get github.com/bborbe/mysql-backup

## Run Backup

One time

```
mysql-backup \
-logtostderr \
-v=2 \
-host=localhost \
-port=5432 \
-lock=/tmp/lock \ 
-username=mysql \
-password=S3CR3T \
-database=db \
-targetdir=/backup \
-name=mysql \
-one-time
```

Cron

```
mysql-backup \
-logtostderr \
-v=2 \
-host=localhost \
-port=5432 \
-lock=/tmp/lock \ 
-username=mysql \
-password=S3CR3T \
-database=db \
-targetdir=/backup \
-name=mysql \
-wait=1h
```
