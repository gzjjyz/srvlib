/**
 * @Author: ChenJunJi
 * @Desc:
 * @Date: 2021/10/8 20:09
 */

package utils

import (
	"runtime"

	"github.com/gzjjyz/trace"
	"github.com/petermattis/goid"

	"github.com/gzjjyz/logger"
)

// ProtectRun 保护方式允许一个函数
func ProtectRun(fn func()) {
	if nil == fn {
		return
	}
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		if nil == err {
			return
		}
		switch err.(type) {
		case runtime.Error: // 运行时错误
			logger.LogStack("runtime error:%v", err)
		default: // 非运行时错误
			logger.LogStack("error:%v", err)
		}
	}()
	fn()
}

// ProtectGo in fact, we could add trace system into there, consider it!!除了recover保护程序以外,这里还可以考虑加入系统链路信息
func ProtectGo(fn func()) {
	var (
		traceId string
		ok      bool
	)
	if traceId, ok = trace.Ctx.GetCurGTrace(goid.Get()); !ok {
		traceId = trace.GenTraceId()
	}
	go func() {
		gid := goid.Get()
		trace.Ctx.SetCurGTrace(gid, traceId)
		defer trace.Ctx.RemoveGTrace(gid)
		ProtectRun(fn)
	}()
}
