package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gzjjyz/logger"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func initScanTempList() []string {
	// 参考 json 化日志结构 {"level":3,"timestamp":"09-20 15:59:58.8011","app_name":"gameserver","content":"The processing time exceeds two milliseconds!! 3.792248ms (2, 64)","trace_id":"UNKNOWN","file":"jjyz/gameserver/logicworker/entity/player.go","line":951,"func":"(*Player).DoNetMsg","prefix":"[srv:100 pf:1]","stack":""}
	var temp = "{\"level\":%d,"
	var scanTempList []string
	for i := logger.ErrorLevel; i <= logger.FatalLevel; i++ {
		scanTempList = append(scanTempList, fmt.Sprintf(temp, i))
	}
	return scanTempList
}

func main() {
	tempList := initScanTempList()
	in := make(chan string)
	out := make(chan string)
	go func() {
		for {
			select {
			case s := <-in:
				for _, temp := range tempList {
					if !strings.HasPrefix(s, temp) {
						continue
					}
					out <- s
					break
				}
			}
		}
	}()

	go func() {
		for {
			content := <-out
			data := map[string]interface{}{
				"msgtype": "text",
				"text": map[string]string{
					"content": content,
				},
			}
			j, err := json.Marshal(data)
			if err != nil {
				fmt.Println(err)
			}
			_, err = http.Post(
				"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=4782b8b5-24c8-4af6-8b25-e81be942c841",
				"application/json", bytes.NewBuffer(j))
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	mapFNameToSeek := make(map[string]int64)

	reOpenFiles := func() {
		var rmFNames []string
		for fName := range mapFNameToSeek {
			if !strings.Contains(fName, time.Now().Format("01-02")) {
				rmFNames = append(rmFNames, fName)
			}
		}

		for _, fname := range rmFNames {
			delete(mapFNameToSeek, fname)
		}

		dir := os.Getenv("TLOGDIR")
		files, err := os.ReadDir(dir)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, file := range files {
			if !strings.Contains(file.Name(), time.Now().Format("01-02")) {
				continue
			}

			if !strings.HasSuffix(file.Name(), ".log") {
				continue
			}

			if _, ok := mapFNameToSeek[file.Name()]; !ok {
				mapFNameToSeek[file.Name()] = 0
			}
		}
	}

	for {
		reOpenFiles()

		for fname, seek := range mapFNameToSeek {
			f, err := os.Open(os.Getenv("TLOGDIR") + fname)
			if err != nil {
				fmt.Println(err)
				continue
			}

			finfo, err := f.Stat()
			if err != nil {
				fmt.Println(err)
				continue
			}

			if seek == 0 {
				seek = finfo.Size()
				mapFNameToSeek[fname] = seek
				f.Close()
				continue
			} else if finfo.Size() <= seek {
				f.Close()
				continue
			}

			_, err = f.Seek(seek, 0)
			if err != nil {
				f.Close()
				fmt.Println(err)
				continue
			}

			rd := bufio.NewReader(f)
			for {
				line, err := rd.ReadString('\n')
				var isEOF bool
				if err != nil {
					if io.EOF != err {
						fmt.Println(err)
					} else {
						isEOF = true
					}
				}

				if len(line) > 0 {
					in <- line
					n := len([]byte(line))
					seek = seek + int64(n)
					mapFNameToSeek[fname] = seek
				}

				if isEOF {
					break
				}
			}
			f.Close()
		}

		time.Sleep(time.Second)
	}
}
