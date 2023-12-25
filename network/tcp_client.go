package network

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gzjjyz/logger"
)

type TCPClient struct {
	sync.Mutex
	Addr            string
	ConnNum         int
	ConnectInterval time.Duration
	PendingWriteNum int
	AutoReconnect   bool
	NewAgent        func(*TCPConn) Agent
	conns           ConnSet
	wg              sync.WaitGroup
	closeFlag       atomic.Bool

	// msg parser
	LenMsgLen    int
	MinMsgLen    uint32
	MaxMsgLen    uint32
	LittleEndian bool
	msgParser    *MsgParser
}

func (client *TCPClient) Start() {
	client.init()

	for i := 0; i < client.ConnNum; i++ {
		client.wg.Add(1)
		go client.connect()
	}
}

func (client *TCPClient) init() {
	client.Lock()
	defer client.Unlock()

	if client.ConnNum <= 0 {
		client.ConnNum = 1
		logger.LogInfo("invalid ConnNum, reset to %v", client.ConnNum)
	}
	if client.ConnectInterval <= 0 {
		client.ConnectInterval = 3 * time.Second
		logger.LogInfo("invalid ConnectInterval, reset to %v", client.ConnectInterval)
	}
	if client.PendingWriteNum <= 0 {
		client.PendingWriteNum = 100
		logger.LogInfo("invalid PendingWriteNum, reset to %v", client.PendingWriteNum)
	}
	if client.NewAgent == nil {
		logger.LogFatal("NewAgent must not be nil")
	}
	if client.conns != nil {
		logger.LogFatal("client is running")
	}

	client.conns = make(ConnSet)
	client.closeFlag.Store(false)

	// msg parser
	msgParser := NewMsgParser()
	msgParser.SetMsgLen(client.LenMsgLen, client.MinMsgLen, client.MaxMsgLen)
	msgParser.SetByteOrder(client.LittleEndian)
	client.msgParser = msgParser
}

func (client *TCPClient) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", client.Addr)
		if err == nil || client.closeFlag.Load() {
			return conn
		}

		logger.LogError("connect to %v error: %v", client.Addr, err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *TCPClient) connect() {
	defer client.wg.Done()

reconnect:
	conn := client.dial()
	if conn == nil {
		return
	}

	client.Lock()
	if client.closeFlag.Load() {
		client.Unlock()
		conn.Close()
		return
	}
	client.conns[conn] = struct{}{}
	client.Unlock()

	tcpConn := newTCPConn(conn, client.PendingWriteNum, client.msgParser)
	agent := client.NewAgent(tcpConn)
	agent.Run()

	// cleanup
	tcpConn.Close()
	client.Lock()
	delete(client.conns, conn)
	client.Unlock()
	agent.OnClose()

	if client.AutoReconnect {
		time.Sleep(client.ConnectInterval)
		goto reconnect
	}
}

func (client *TCPClient) Close() {
	client.Lock()
	client.closeFlag.Store(true)
	for conn := range client.conns {
		conn.Close()
	}
	client.conns = nil
	client.Unlock()

	client.wg.Wait()
}
