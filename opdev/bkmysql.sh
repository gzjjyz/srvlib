#!/bin/bash
HOST=0.0.0.0                             #ip地址
USER=root                                #数据库用户名
PASSWORD=game@2023                       #数据库密码
DATABASE=bk_test                         #数据库名字
BACKUP_PATH=/tmp/bkdata/${DATABASE}      #备份路径
logfile=/tmp/bkdata/${DATABASE}/data.log #日志路径
DATE=$(date '+%Y%m%d')                   #当前日期
LastDay=7                                #七天前的都备份

#如果不存在该目录，就创建
stat $BACKUP_PATH
if [ $? != 0 ]; then
  mkdir -p "$BACKUP_PATH"
  echo "create directory $BACKUP_PATH"
fi

# 备份数据库
bkMysqlTable() {

  # 拿到这几天的日期
  lastDayList=()
  for ((i = 0; i <= ${LastDay}; i++)); do
    lastDayList[$i]=$(date -d "-$i day" +%Y%m%d)
  done

  #进入到备份目录
  cd $BACKUP_PATH

  #遍历数据库中的数据表
  for table in $DATABASE; do
    echo "start backup database $DATABASE" >>$logfile

    #获取表名
    table=$(mysql -h $HOST -u $USER -p$PASSWORD $DATABASE -e "show tables;" | sed '1d')
    for tb in $table; do
      echo "start backup table ${tb}" >>$logfile

      # 是否保留 - 按时间去判断
      reserve=""
      for day in ${lastDayList[*]}; do
        if [[ $tb =~ $day ]]; then
          reserve="true"
          echo "reserve table ${tb}" >>$logfile
          break
        else
          continue
        fi
      done

      if [[ $reserve != "" ]]; then
        continue
      fi

      # 开始备份
      echo "backup table ${tb}" >>$logfile

      # 备份文件名
      DUMPNAME=""$DATE"_bk_"$tb".sql"

      # 备份
      mysqldump -h $HOST -u $USER -p$PASSWORD $DATABASE $tb >$DUMPNAME

      # 备份成功
      if [ $? = 0 ]; then
        echo "$DUMPNAME backup Successful!" >>$logfile
        echo "start drop table ${tb}" >>$logfile
        mysql -h $HOST -u $USER -p$PASSWORD $DATABASE -e "drop table if exists $tb"
        echo "end drop table ${tb}" >>$logfile
      else
        echo "$DUMPNAME backup fail!" >>$logfile
      fi
      echo "end backup table ${tb}" >>$logfile
    done
    echo "end backup database $DATABASE" >>$logfile
  done
}

# 压缩文件
tarBkDir() {
  cd /tmp
  tar -zcvf "bk"_${DATE}"_"${DATABASE}.tar.gz ${BACKUP_PATH}
}

# 导入
# $1 sql文件
importSql() {
  mysql -h $HOST -u $USER -p$PASSWORD $DATABASE < $1
}

bkMysqlTable
tarBkDir

# 清理
rm -rf ${BACKUP_PATH}
