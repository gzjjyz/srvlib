package v2

import (
	"fmt"
	"testing"
	"time"

	"github.com/gzjjyz/logger"
)

func TestWorker(t *testing.T) {
	logger.InitLogger(logger.WithAppName("test"))
	worker := NewWorker(100, func() {
		fmt.Println("abc")
	})

	go worker.GoStart()

	worker.RegisterMsgHandler(123, func(param ...interface{}) {
		logger.LogDebug("i consume msg 123.%d\n", param[0])
	})

	worker.RegisterMsgHandler(234, func(param ...interface{}) {
		logger.LogDebug("i consume msg 234.%d\n", param[0])
	})

	now := time.Now()
	i := 0
	for {
		if i >= 100 {
			break
		}
		i++
		worker.SendMsg(123, i)
		worker.SendMsg(234, i)
	}

	logger.LogDebug("%d", i)
	for {
		if i >= 200 {
			break
		}
		i++
		worker.SendMsg(234, i)
	}

	for {
		if i >= 300 {
			break
		}
		i++
		worker.SendMsg(123, i)
	}

	logger.LogDebug(time.Since(now).String())
	worker.Stop()
}
