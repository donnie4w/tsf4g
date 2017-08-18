/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package register

import (
	"bytes"
	//	"errors"
	"fmt"
	"time"

	. "github.com/donnie4w/tsf/packet"
	. "github.com/donnie4w/tsf/utils"
)

type TsFTIOStreamTransport struct {
	timeout   time.Duration
	outBs     []byte
	inBs      []byte
	buf       *bytes.Buffer
	inBuf     *bytes.Buffer
	regist    *Regist
	serviceid string
}

func NewTsFTIOStreamTransport(regist *Regist, serviceid string) (tt *TsFTIOStreamTransport, err error) {
	tt = &TsFTIOStreamTransport{buf: bytes.NewBuffer([]byte{}), regist: regist, serviceid: serviceid}
	return
}

// Sets the socket timeout
func (p *TsFTIOStreamTransport) SetTimeout(timeout time.Duration) error {
	p.timeout = timeout
	return nil
}

// Connects the socket, creating a new socket object if necessary.
func (p *TsFTIOStreamTransport) Open() error {
	//	if p.IsOpen() {
	//		return errors.New("Socket already connected.")
	//	}
	return nil
}

// Returns true if the connection is open
func (p *TsFTIOStreamTransport) IsOpen() bool {
	//	if p.regist == nil {
	//		return false
	//	}
	return p.regist.IsOpen()
}

// Closes the socket.
func (p *TsFTIOStreamTransport) Close() error {
	return p.regist.Close()
}

func (p *TsFTIOStreamTransport) Read(buf []byte) (i int, err error) {
	if p.inBuf == nil {
		p.inBuf = bytes.NewBuffer([]byte{})
		pk := NewPacket()
		pk.PType = Register
		pk.Body = p.buf.Bytes()
		pk.SeqId = time.Now().UnixNano()
		pk.Serviceid = Md5(p.serviceid)
		bs, err := p.regist.Register(pk)
		if err == nil {
			fmt.Println("bs==>", bs)
			p.inBuf.Write(bs)
		}
	}
	p.buf = bytes.NewBuffer([]byte{})
	return p.inBuf.Read(buf)
}

func (p *TsFTIOStreamTransport) Write(buf []byte) (i int, err error) {
	p.inBuf = nil
	return p.buf.Write(buf)
}

func (p *TsFTIOStreamTransport) Flush() error {
	return nil
}

func (p *TsFTIOStreamTransport) RemainingBytes() (num_bytes uint64) {
	const maxSize = ^uint64(0)
	return maxSize // the thruth is, we just don't know unless framed is used
}
