/**
 * Copyright 2017 tsf Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package tsocket

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/donnie4w/go-logger/logger"
	. "github.com/donnie4w/tsf/packet"
	. "github.com/donnie4w/tsf/utils"
)

type Tsocket struct {
	Host            string
	Port            int
	I_              int32
	IsClosed        bool
	conn            *net.TCPConn
	ConnectTimeout_ int
	SocketTimeout_  int
}

func NewTsocket(host string, port int) (socket *Tsocket, err error) {
	socket = &Tsocket{Host: host, Port: port, I_: 0, ConnectTimeout_: 3, SocketTimeout_: 3}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprint(host, ":", port))
	if err != nil {
		return nil, err
	}
	socket.conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return
}

func (this *Tsocket) ReadPacket() (p *Packet, e error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			e = errors.New(fmt.Sprint(err))
		}
	}()
	headBs := make([]byte, 4)
	n, err := this.conn.Read(headBs)
	if err != nil {
		return nil, err
	}
	if n != 4 {
		return nil, errors.New("readpacket err")
	}
	length := int(BytesToInt(headBs))
	buf := bytes.NewBuffer([]byte{})
	i := 0
	for {
		bs := make([]byte, length-buf.Len())
		i, err = this.conn.Read(bs)
		if err == nil {
			buf.Write(bs[0:i])
		} else {
			panic(err.Error())
		}
		if length == buf.Len() {
			break
		}
	}
	bs := buf.Bytes()
	fmt.Println("2bs====>", bs)
	p = Wrap(bs)
	return
}

func (this *Tsocket) WritePacket(p *Packet) (err error) {
	_, err = this.conn.Write(p.GetPacket())
	return
}

func (this *Tsocket) Close() error {
	return this.conn.Close()
}
