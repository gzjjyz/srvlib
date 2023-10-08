#!/bin/bash
SaveLogDir=/data/log/gzjjyz    #日志存放
DATE=$(date '+%Y%m%d')         #当前日期
BACKUP_PATH=/tmp/bklog/${DATE} #备份路径
LastDay=7                      #七天前的都备份

#如果不存在该目录，就创建
stat $BACKUP_PATH
if [ $? != 0 ]; then
  mkdir -p "$BACKUP_PATH"
  echo "create directory $BACKUP_PATH"
fi

# 备份日志
bkLog() {

  # 拿到这几天的日期
  lastDayList=()
  for ((i = 0; i <= ${LastDay}; i++)); do
    lastDayList[$i]=$(date -d "-$i day" +%m-%d)
  done

  #进入到备份目录
  cd $SaveLogDir

  echo "start backup log $SaveLogDir"

  filename=$(ls)
  for fn in $filename; do
    # 是否保留 - 按时间去判断
    reserve=""
    for day in ${lastDayList[*]}; do
      if [[ $fn =~ $day ]]; then
        reserve="true"
        break
      else
        continue
      fi
    done

    if [[ $reserve != "" ]]; then
      continue
    fi

    # 开始备份
    echo "backup log file ${fn}"

    # 迁移到某个文件夹下
    mv -f ${fn} ${BACKUP_PATH}
  done

  # 打包
  cd ${BACKUP_PATH}/.. && tar -zcvf "bklog"${DATE}.tar.gz ${DATE}

  rm -rf ${BACKUP_PATH}

  echo "end backup log $SaveLogDir"

}

bkLog
