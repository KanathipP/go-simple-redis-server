package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/tidwall/resp"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
	delCh chan *Peer
}

func (p *Peer) Send(msg []byte) (int, error) {
	return p.conn.Write(msg)
}

func NewPeer(conn net.Conn, msgCh chan Message, delCh chan *Peer) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
		delCh: delCh,
	}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)

	for {
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			p.delCh <- p
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			arr := v.Array()
			commandName := arr[0].String()
			var cmd Command

			switch commandName {
			case CommandClient:
				cmd = ClientCommand{
					val: arr[1].String(),
				}
			case CommandGET:
				if len(arr) != 2 {
					return fmt.Errorf("invalid number of variables for GET command")
				}
				cmd = GetCommand{
					key: arr[1].Bytes(),
				}
			case CommandSET:
				if len(arr) != 3 {
					return fmt.Errorf("invalid number of variables for SET command")
				}
				cmd = SetCommand{
					key: arr[1].Bytes(),
					val: arr[2].Bytes(),
				}
			case CommandHELLO:
				cmd = HelloCommand{
					val: arr[1].String(),
				}
			default:
				fmt.Printf("got unknown command => %+v\n", v.Array())
			}
			p.msgCh <- Message{
				cmd:  cmd,
				peer: p,
			}

		}
	}
	return nil
}
