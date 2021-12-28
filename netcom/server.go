package netcom

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ruraomsk/TLServer/logger"
	"github.com/ruraomsk/brams/drive"
)

var remotes struct {
	sync.Mutex
	regs map[string]*remote
}

type remote struct {
	socket net.Conn
	db     *drive.Db
}

func ServerCommections(port int, timeout time.Duration, stop chan interface{}) {
	remotes.regs = make(map[string]*remote)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Error.Printf("Ошибка открытия порта %d %s", port, err.Error())
		stop <- 1
	}
	for {
		socket, err := ln.Accept()
		if err != nil {
			logger.Error.Printf("Ошибка accept %s", err.Error())
			continue
		}
		go workerConnect(socket, timeout)
	}
}
func workerConnect(socket net.Conn, timeout time.Duration) {
	registryConn(socket)
	defer closeConnect(socket)
	reader := bufio.NewReader(socket)
	for {
		socket.SetReadDeadline(time.Now().Add(timeout))
		c, err := reader.ReadString('\n')
		if err != nil {
			logger.Error.Printf("При чтении команды от %s %s", socket.RemoteAddr().String(), err.Error())

		}
		if c[0:1] == "0" {
			//Это keepalive
			continue
		}

	}

}

func registryConn(socket net.Conn) {
	if _, is := remotes.regs[socket.LocalAddr().String()]; !is {
		logger.Error.Printf("Remote %s уже зарегистрирован!", socket.LocalAddr().String())
		return
	}
	r := new(remote)
	r.socket = socket
	r.db = nil

}
func closeConnect(socket net.Conn) {
	remotes.Lock()
	r, is := remotes.regs[socket.LocalAddr().String()]
	if !is {
		return
	}
	r.socket.Close()
	if r.db != nil {
		r.db.Close()
	}
	delete(remotes.regs, socket.LocalAddr().String())
}
