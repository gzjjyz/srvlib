#!/bin/bash
GTLOGDIR=/data/log/gzjjyz # 重新拿全局日志路径变量
HOST=192.168.61.231
USER=server
SSH_PARENT_DIR=/home/server             # scp 上传的目录
SSH_SRV_DIR=${SSH_PARENT_DIR}/bkcronbiz # 解压后文件的所在路径

mkdir -p bkcronbiz

sh ./gobuild.sh app ./bkcron

mv app bkcronbiz/
cp -rf ./bkcron/cfg.yaml bkcronbiz/
cp -rf ./bkcron/bklog.sh bkcronbiz/
cp -rf ./bkcron/bkmysql.sh bkcronbiz/

echo "
[program:bkcron] ;
directory =${SSH_SRV_DIR}
command =${SSH_SRV_DIR}/app -p ${SSH_SRV_DIR}
autostart = true
autorestart = true
startsecs = 10
startretries = 3
user = server
redirect_stderr = true
stdout_logfile=/var/log/bkcron_std.log
stderr_logfile=/var/log/bkcron_err.log
stdout_logfile_maxbytes = 20MB
stdout_logfile_backups = 20
environment=TLOGDIR=${GTLOGDIR}
" >bkcronbiz/bkcron.conf

tar -zcvf bkcronbiz.tar.gz bkcronbiz

ssh ${USER}@${HOST} "
  mkdir -p ${SSH_SRV_DIR}
"

scp -rp bkcronbiz.tar.gz server@${HOST}:${SSH_PARENT_DIR}

ssh ${USER}@${HOST} "cd ${SSH_PARENT_DIR} && tar -zxvf bkcronbiz.tar.gz ;
  chmod +x -R bkcronbiz;
  cp -rf ${SSH_SRV_DIR}/bkcron.conf /etc/supervisord.d ;
  supervisorctl update
  supervisorctl restart bkcron
  rm -rf bkcronbiz.tar.gz
"

rm -rf ./bkcronbiz
rm -rf ./bkcronbiz.tar.gz
