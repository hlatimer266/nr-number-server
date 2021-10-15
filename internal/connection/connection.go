package connection

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/hlatimer266/nr-number-server/internal/cache"
	"github.com/hlatimer266/nr-number-server/internal/write"
)

const (
	maxConnections = 5
	shutDownMsg    = "terminate$"
)

type Listener struct {
	Ctx          context.Context
	Port         string
	Listen       *net.TCPListener
	ClientWG     sync.WaitGroup
	AllWG        sync.WaitGroup
	NumConns     int
	NumCache     *cache.NumberCache
	StopWriter   context.CancelFunc
	StopReport   context.CancelFunc
	StopListener context.CancelFunc
	ConnCount    NumberConnections
}

type NumberConnections struct {
	Count int
	sync.RWMutex
}

func NewListener(p string) Listener {
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", p)
	l, _ := net.ListenTCP("tcp", tcpAddr)
	return Listener{
		Port:      p,
		NumConns:  0,
		ClientWG:  sync.WaitGroup{},
		Listen:    l,
		ConnCount: NumberConnections{Count: 0},
	}
}

func (l *Listener) readMessage(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		message = message[0 : len(message)-1] // Remove the new line delimiter
		msgStr := string(message)

		if msgStr == shutDownMsg {
			fmt.Println("Shutting down server...")
			l.StopAll() // gracefully shut down server upon receiving termination request
			break
		}
		// add message processing to the wait group + launch new thread to handle message
		l.ClientWG.Add(1)
		go func() {
			if err := write.NumCache(msgStr, l.NumCache); err != nil {
				fmt.Println(err.Error())
				if _, ok := err.(*os.PathError); ok {
					l.StopAll() // shut down server if file can't be created / opened
				}
				conn.Close() // close this connection if bad input is received
			}
			l.ClientWG.Done()
		}()

	}
}

func (l *Listener) ManageClient() error {
	for {
		conn, err := l.Listen.Accept() // wait and accept a new connection
		if err != nil {
			return err
		}

		l.ConnCount.Add()
		if l.ConnCount.Count < maxConnections {
			go l.readMessage(conn) // process messages from new connetion
		}
		l.ConnCount.Subtract()
	}
}

func (c *NumberConnections) Add() {
	c.Lock()
	c.Count++
	c.Unlock()
}

func (c *NumberConnections) Subtract() {
	c.Lock()
	c.Count--
	c.Unlock()
}

func (l *Listener) StopAll() {
	l.StopReport()
	l.StopWriter()
	l.StopListener()
}

func (l *Listener) WaitCtxFinish() {
	<-l.Ctx.Done()
	l.StopAll()
}
