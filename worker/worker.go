/**
 * @Author: ChenJunJi
 * @Desc:
 * @Date: 2021/8/27 16:23
 */

package worker

import (
	"fmt"
	"github.com/gzjjyz/logger"
	"github.com/gzjjyz/srvlib/alg/queue_list"
	"github.com/gzjjyz/srvlib/utils"
	"log"
	"sync/atomic"
	"time"
)

type MsgHdlType func(param ...interface{})

const SleepTime = time.Millisecond * 10

type MsgSt struct {
	MsgId uint32        // 消息id
	Param []interface{} // 参数
}

func (m *MsgSt) String() string {
	str := fmt.Sprintf("id:%d ", m.MsgId)
	for i, p := range m.Param {
		strSeg := fmt.Sprintf("param%d:%+v ", i+1, p)
		if len(strSeg) > 20 {
			strSeg = strSeg[:20] + "..."
		}
		str += strSeg
	}
	return str
}

type Worker struct {
	sleep time.Duration

	stop     int32
	exit_    chan bool
	loopFunc func()
	msgList  *queue_list.QueueListSt
	mHdl     map[uint32]MsgHdlType

	stopRecvFromGate int32
	exitGate_        chan bool
}

func NewWorker(msgCapacity uint32, loopFunc func()) *Worker {
	if nil == loopFunc {
		log.Fatalf("Worker Loop Func IsNil")
	}
	worker := &Worker{}
	worker.exit_ = make(chan bool)
	worker.exitGate_ = make(chan bool)
	worker.loopFunc = loopFunc
	worker.msgList = queue_list.NewQueueList(msgCapacity)
	worker.mHdl = make(map[uint32]MsgHdlType)
	worker.sleep = SleepTime
	return worker
}

func (worker *Worker) RegisterMsgHandler(msgId uint32, hdl MsgHdlType) {
	if nil == hdl {
		log.Fatalf("注册消息处理函数为空, 消息id=%d", msgId)
		return
	}
	_, repeat := worker.mHdl[msgId]
	if repeat {
		log.Fatalf("注册消息重复, 消息id=%d", msgId)
		return
	}
	worker.mHdl[msgId] = hdl
}

func (worker *Worker) SendMsg(id uint32, params ...interface{}) {
	if atomic.LoadInt32(&worker.stop) == 1 {
		return
	}
	st := &MsgSt{
		MsgId: id,
		Param: params,
	}
	worker.msgList.Append(st)
}

func (worker *Worker) SendMsgFromGate(id uint32, params ...interface{}) {
	if atomic.LoadInt32(&worker.stopRecvFromGate) == 1 {
		return
	}
	st := &MsgSt{
		MsgId: id,
		Param: params,
	}
	worker.msgList.Append(st)
}

func (worker *Worker) SetSleep(t time.Duration) {
	worker.sleep = t
}

func (worker *Worker) GoStart() bool {
	utils.ProtectGo(func() {
	out:
		for {
			select {
			case <-worker.exit_:
				atomic.StoreInt32(&worker.stop, 1)
				break out
			case <-worker.exitGate_:
				atomic.StoreInt32(&worker.stopRecvFromGate, 1)
				break out
			default:
				worker.loop()
			}
		}
	})
	return true
}

func (worker *Worker) Stop() {
	worker.exit_ <- true
	worker.ProcessMsg()
}

func (worker *Worker) StopGate() {
	worker.exitGate_ <- true
	worker.ProcessMsg()
}

func (worker *Worker) loop() {
	defer func() {
		if err := recover(); err != nil {
			logger.Stack("循环中出现错误:%v", err)
		}
	}()

	utils.ProtectRun(worker.loopFunc)
	worker.ProcessMsg()

	if worker.sleep > 0 {
		time.Sleep(worker.sleep)
	}
}

func (worker *Worker) ProcessMsg() {
	worker.msgList.Flush()
	worker.msgList.Traverse(func(args interface{}) {
		if msg, ok := args.(*MsgSt); ok {
			t := time.Now()
			if fn, ok := worker.mHdl[msg.MsgId]; ok {
				utils.ProtectRun(func() {
					fn(msg.Param[:]...)
				})
			}
			if since := time.Since(t); since > 20*time.Millisecond {
				logger.Debug("process msg end! id:%d, cost:%v", msg.MsgId, since)
			}
		}
	})
}
