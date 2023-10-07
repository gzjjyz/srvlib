package main

import (
	"fmt"
	"github.com/gzjjyz/logger"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"os"
	"path"
	"strings"
	"testing"
)

func Test_upload(t *testing.T) {
	logger.InitLogger()
	err := loadGlobal()
	if err != nil {
		logger.Errorf("err:%v", err)
		return
	}
	obsClient, err = obs.New(global.Obs.AccessKey, global.Obs.SecretKey, global.Obs.EndPoint)
	if err != nil {
		logger.Errorf("err:%v", err)
		return
	}

	vo := global.BkList[1]
	files, err := os.ReadDir(vo.DirPath)
	if err != nil {
		logger.Errorf("err:%v", err)
		return
	}

	var backupLogs []os.DirEntry
	for _, file := range files {
		if vo.Prefix == "" && vo.Suffix == "" {
			continue
		}
		if !strings.HasPrefix(file.Name(), vo.Prefix) || !strings.HasSuffix(file.Name(), vo.Suffix) {
			continue
		}
		backupLogs = append(backupLogs, file)
	}

	if len(backupLogs) == 0 {
		logger.Info("not found backup file")
		return
	}

	// 上传
	for _, file := range backupLogs {
		err := upload(global.Obs.Bucket, file.Name(), path.Join(vo.DirPath, file.Name()))
		if err != nil {
			logger.Errorf("err:%v", err)
			continue
		}

		if vo.AfterUploadRm {
			if err = os.Remove(path.Join(vo.DirPath, file.Name())); err != nil {
				logger.Errorf("del %s file fail,err:%v", path.Join(vo.DirPath, file.Name()), err)
			}
		}
	}
}

func Test_bucket(t *testing.T) {
	logger.InitLogger()
	err := loadGlobal()
	if err != nil {
		logger.Errorf("err:%v", err)
		return
	}
	obsClient, err := obs.New(global.Obs.AccessKey, global.Obs.SecretKey, global.Obs.EndPoint)
	if err != nil {
		fmt.Printf("Create obsClient error, errMsg: %s", err.Error())
	}

	var bucketname = "gz-xianxia-v2-file"

	// 列举桶列表
	output, err := obsClient.HeadBucket(bucketname)
	if err == nil {
		fmt.Printf("Head bucket(%s) successful!\n", bucketname)
		fmt.Printf("RequestId:%s\n", output.RequestId)
		return
	}
	fmt.Printf("Head bucket(%s) fail!\n", bucketname)
	if obsError, ok := err.(obs.ObsError); ok {
		fmt.Println("An ObsError was found, which means your request sent to OBS was rejected with an error response.")
		fmt.Println(obsError.Error())
	} else {
		fmt.Println("An Exception was found, which means the client encountered an internal problem when attempting to communicate with OBS, for example, the client was unable to access the network.")
		fmt.Println(err)
	}
}
