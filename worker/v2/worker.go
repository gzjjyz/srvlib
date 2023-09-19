package v2

import (
	"github.com/gzjjyz/logger"
	"github.com/gzjjyz/srvlib/trace"
	"github.com/gzjjyz/srvlib/utils"
	work "github.com/gzjjyz/srvlib/worker"
	"github.com/petermattis/goid"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

const (
	revBatchMsgMaxWait    = time.Millisecond * 10
	loopEventProcInterval = time.Millisecond * 10
)

type TracedMsg struct {
	TraceId string
	*work.MsgSt
}

type Worker struct {
	stopped             atomic.Bool
	exitCh              chan bool
	exitWait            sync.WaitGroup
	loopFunc            func()
	mHdl                map[uint32]work.MsgHdlType
	msgCh               chan *TracedMsg
	procBatchMsgMaxSize uint32
	firstMsgInLoop      *TracedMsg
	stoppedRecvFromGate atomic.Bool
	exitGateCh          chan bool
}

func NewWorker(msgCapacity uint32, loopFunc func()) *Worker {
	if nil == loopFunc {
		log.Fatalf("Worker Loop Func IsNil")
	}
	worker := &Worker{}
	worker.exitCh = make(chan bool)
	worker.exitGateCh = make(chan bool)
	worker.loopFunc = loopFunc
	worker.msgCh = make(chan *TracedMsg, msgCapacity)
	worker.procBatchMsgMaxSize = msgCapacity
	worker.mHdl = make(map[uint32]work.MsgHdlType)
	return worker
}

func (w *Worker) RegisterMsgHandler(msgId uint32, hdl work.MsgHdlType) {
	if nil == hdl {
		log.Fatalf("注册消息处理函数为空, 消息id=%d", msgId)
		return
	}
	_, repeat := w.mHdl[msgId]
	if repeat {
		log.Fatalf("注册消息重复, 消息id=%d", msgId)
		return
	}
	w.mHdl[msgId] = hdl
}

func (w *Worker) doSendMsg(id uint32, params ...interface{}) {
	var (
		traceId string
		ok      bool
	)
	if traceId, ok = trace.Ctx.GetCurGTrace(goid.Get()); !ok {
		traceId = trace.GenTraceId()
	}
	st := &TracedMsg{
		MsgSt: &work.MsgSt{
			MsgId: id,
			Param: params,
		},
		TraceId: traceId,
	}
	w.msgCh <- st
}

func (w *Worker) SendMsg(id uint32, params ...interface{}) {
	if w.stopped.Load() {
		return
	}
	w.doSendMsg(id, params...)
}

func (w *Worker) SendMsgFromGate(id uint32, params ...interface{}) {
	if w.stoppedRecvFromGate.Load() {
		return
	}
	w.doSendMsg(id, params...)
}

func (w *Worker) GoStart() bool {
	utils.ProtectGo(func() {
		defer trace.Ctx.RemoveGTrace(goid.Get())

		doLoopFuncTk := time.NewTicker(loopEventProcInterval)
		defer doLoopFuncTk.Stop()
	out:
		for {
			select {
			case <-w.exitCh:
				w.stopped.Store(true)
				w.ProcessMsg(w.FetchAndMergeBatch(nil))
				w.exitWait.Done()
				break out
			case <-w.exitGateCh:
				w.stoppedRecvFromGate.Store(true)
				w.ProcessMsg(w.FetchAndMergeBatch(nil))
				w.exitWait.Done()
				break out
			case w.firstMsgInLoop = <-w.msgCh:
				w.loop()
				// let tick restart to calc interval
				doLoopFuncTk.Stop()
				doLoopFuncTk = time.NewTicker(loopEventProcInterval)
			case <-doLoopFuncTk.C:
				w.loop()
			}
		}
	})
	return true
}

func (w *Worker) Stop() {
	w.exitWait.Add(1)
	w.exitCh <- true
	w.exitWait.Wait()
}

func (w *Worker) StopGate() {
	w.exitWait.Add(1)
	w.exitGateCh <- true
	w.exitWait.Wait()
}

func (w *Worker) loop() {
	defer func() {
		if err := recover(); err != nil {
			logger.Stack("循环中出现错误:%v", err)
		}
	}()

	var msgList []*TracedMsg
	if w.firstMsgInLoop != nil {
		msgList = w.FetchAndMergeBatch([]*TracedMsg{w.firstMsgInLoop})
		w.firstMsgInLoop = nil
	} else {
		msgList = w.FetchAndMergeBatch(nil)
	}
	utils.ProtectRun(w.loopFunc)
	w.ProcessMsg(msgList)
}

func (w *Worker) ProcessMsg(msgList []*TracedMsg) {
	gid := goid.Get()
	for _, msg := range msgList {
		t := time.Now()
		if fn, ok := w.mHdl[msg.MsgId]; ok {
			trace.Ctx.SetCurGTrace(gid, msg.TraceId)
			utils.ProtectRun(func() {
				fn(msg.Param[:]...)
			})
		}
		if since := time.Since(t); since > 20*time.Millisecond {
			logger.Debug("process msg end! id:%d, cost:%v", msg.MsgId, since)
		}
	}
}

func (w *Worker) FetchAndMergeBatch(msgList []*TracedMsg) []*TracedMsg {
	t := time.Now()
	for {
		select {
		case msg := <-w.msgCh:
			msgList = append(msgList, msg)
			if uint32(len(msgList)) >= w.procBatchMsgMaxSize {
				goto out
			}
			if since := time.Since(t); since > revBatchMsgMaxWait {
				goto out
			}
		default:
			goto out
		}
	}
out:
	return msgList
}
