package register

import (
	"fmt"

	. "github.com/donnie4w/tsf/packet"
	. "github.com/donnie4w/tsf4g/tsocket"
)

type Regist struct {
	_sokcet *Tsocket
}

func NewRegist(host string, port int) (r *Regist, err error) {
	var socket *Tsocket
	socket, err = NewTsocket(host, port)
	r = &Regist{_sokcet: socket}
	return
}

func (this *Regist) RegisterAuth(pwd string) (err error) {
	p := NewPacket()
	p.PType = Auth
	p.Body = []byte(pwd)
	err = this._sokcet.WritePacket(p)
	readpacket, err := this._sokcet.ReadPacket()
	if err == nil {
		fmt.Println(readpacket.ToString())
	}
	return
}

func (this *Regist) Register(p *Packet) (bs []byte, err error) {
	this._sokcet.WritePacket(p)
	b := true
	for b {
		readpacket, err := this._sokcet.ReadPacket()
		if err != nil {
			fmt.Println("p===>", readpacket.ToString())
			return nil, err
		}
		switch readpacket.PType {
		case Ping:
			pk := NewPacket()
			pk.PType = AckPing
			this._sokcet.WritePacket(pk)
		case AckRegister:
			bs = readpacket.Body
			b = false
		default:
		}
	}
	return
}

func (this *Regist) Close() error {
	return this._sokcet.Close()
}

func (this *Regist) IsOpen() bool {
	return !this._sokcet.IsClosed
}
