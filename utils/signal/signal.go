package signal

import (
	"os"
	"os/signal"
	"syscall"
)

var signalChan = make(chan os.Signal, 1)

func init() {
	list := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	}
	signal.Notify(signalChan, list...)
}

func SignalChan() <-chan os.Signal {
	return signalChan
}
