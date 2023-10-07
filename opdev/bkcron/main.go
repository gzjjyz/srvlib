package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/995933447/std-go/scan"
	"github.com/gzjjyz/srvlib/utils/signal"
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	defaultPrefix = "cfg"
	defaultSuffix = "yaml"
)

var obsClient *obs.ObsClient
var global *Cfg

func main() {
	err := loadGlobal()
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	obsClient, err = obs.New(global.Obs.AccessKey, global.Obs.SecretKey, global.Obs.EndPoint)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	crontab := cron.New()
	for _, bk := range global.BkList {
		entryID, err := crontab.AddFunc(bk.Spec, newCronFunc(bk))
		if err != nil {
			fmt.Printf("err:%v\n", err)
			return
		}
		fmt.Printf("start add cron success , entry id is %v\n", entryID)
	}
	crontab.Start()

	signals := signal.SignalChan()
	s := <-signals
	fmt.Printf("received signal [%v]\n", s)
	crontab.Stop()
}

func newCronFunc(vo *CmdVo) func() {
	return func() {
		cmd := exec.Command(vo.Command)
		fmt.Printf("exec: %s\n", cmd.String())

		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err := cmd.Run()
		outStr := outBuf.String()
		errStr := errBuf.String()
		if err != nil {
			fmt.Printf("Err: %s \nStdout: %s \nStderr: %s\n", err, outStr, errStr)
			return
		}
		if outStr != "" {
			fmt.Printf("out:%s\n", outStr)
		}
		if errStr != "" {
			fmt.Printf("err:%s\n", errStr)
		}

		files, err := os.ReadDir(vo.DirPath)
		if err != nil {
			fmt.Printf("err:%v\n", err)
			return
		}

		var backupLogs []os.DirEntry
		for _, file := range files {
			if vo.Prefix == "" && vo.Suffix == "" {
				continue
			}
			fmt.Printf("%s\n ,prefix is %s, suffix is %s \n", file.Name(), vo.Prefix, vo.Suffix)
			if !strings.HasPrefix(file.Name(), vo.Prefix) || !strings.HasSuffix(file.Name(), vo.Suffix) {
				continue
			}
			backupLogs = append(backupLogs, file)
		}

		if len(backupLogs) == 0 {
			fmt.Printf("not found backup file\n")
			return
		}

		// 上传
		if !vo.UploadObs {
			fmt.Printf("not upload")
			return
		}
		for _, file := range backupLogs {
			err := upload(global.Obs.Bucket, file.Name(), path.Join(vo.DirPath, file.Name()))
			if err != nil {
				fmt.Printf("err:%v\n", err)
				continue
			}

			if vo.AfterUploadRm {
				if err = os.Remove(path.Join(vo.DirPath, file.Name())); err != nil {
					fmt.Printf("del %s file fail,err:%v\n", path.Join(vo.DirPath, file.Name()), err)
				}
			}
		}

	}
}

func loadGlobal() error {
	viper.SetConfigName(defaultPrefix)
	viper.SetConfigType(defaultSuffix)
	viper.AddConfigPath(scan.OptStrDefault("p", defaultPath))
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	val := viper.Get("cfg")
	var cfg Cfg
	err = json2St(val, &cfg)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return err
	}
	global = &cfg
	return nil
}

func upload(bucket, objKey, filepath string) error {
	input := &obs.UploadFileInput{}
	input.Bucket = bucket
	input.Key = path.Join(global.Obs.Dir, objKey)
	input.UploadFile = filepath
	input.EnableCheckpoint = true
	input.PartSize = 1024
	input.TaskNum = 5
	output, err := obsClient.UploadFile(input)
	if err == nil {
		fmt.Printf("Upload file(%s) under the bucket(%s) successful!\n", input.UploadFile, input.Bucket)
		fmt.Printf("ETag:%s\n", output.ETag)
		return nil
	}
	fmt.Printf("Upload file(%s) under the bucket(%s) fail!\n", input.UploadFile, input.Bucket)
	if obsError, ok := err.(obs.ObsError); ok {
		fmt.Printf("An ObsError was found, which means your request sent to OBS was rejected with an error response.\n")
		fmt.Printf("err:%v\n", obsError.Error())
	} else {
		fmt.Printf("An Exception was found, which means the client encountered an internal problem when attempting to communicate with OBS, for example, the client was unable to access the network.\n")
		fmt.Printf("err:%v\n", err)
	}
	return nil
}

type Cfg struct {
	BkList   []*CmdVo `json:"bk_list"`
	Obs      *Obs     `json:"obs"`
	PartSize int64    `json:"part_size"`
}

type Obs struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	EndPoint  string `json:"end_point"`
	Bucket    string `json:"bucket"`
	Dir       string `json:"dir"`
}

type CmdVo struct {
	Spec    string `json:"spec"`    // 定时
	Command string `json:"command"` // 执行脚本

	UploadObs     bool   `json:"upload_obs"` // 是否上传 obs
	DirPath       string `json:"dir_path"`   // 指定路径
	Prefix        string `json:"prefix"`     // 指定前缀
	Suffix        string `json:"suffix"`
	AfterUploadRm bool   `json:"after_upload_rm"` // 是否上传完删除文件
}

func json2St(re interface{}, out interface{}) error {
	marshal, err := json.Marshal(re)
	if err != nil {
		fmt.Printf("err is : %v\n", err)
		return err
	}

	err = json.Unmarshal(marshal, out)
	if err != nil {
		fmt.Printf("err is : %v\n", err)
		return err
	}
	return nil
}
