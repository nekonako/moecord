package websocket

import (
	"fmt"
	"net"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/nekonako/moecord/config"
	"github.com/rs/zerolog/log"
)

type conn struct {
	*sync.RWMutex
	conn map[net.Conn]bool
}

var (
	connOnce sync.Once
	connMap  *conn
)

type websocket struct {
	u ws.Upgrader
	l net.Listener
}

func New(c *config.Config) (*websocket, error) {

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.Websocket.Host, c.Websocket.Port))
	if err != nil {
		return nil, err
	}

	u := ws.Upgrader{
		OnHeader: func(key, value []byte) error {
			return err
		},
	}

	connOnce.Do(func() {
		connMap = &conn{
			conn:    map[net.Conn]bool{},
			RWMutex: &sync.RWMutex{},
		}
	})

	return &websocket{
		l: l,
		u: u,
	}, nil

}

func (w *websocket) ListenConnection() {
	log.Info().Msg("listen websocket connection")
	for {
		conn, err := w.l.Accept()
		if err != nil {
			log.Error().Msg(err.Error())
		}
		_, err = w.u.Upgrade(conn)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		connMap.Lock()
		connMap.conn[conn] = true
		connMap.Unlock()
		go w.handleConnection(conn)
	}
}

func (w *websocket) handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		msg, op, err := wsutil.ReadClientData(conn)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
		err = wsutil.WriteServerMessage(conn, op, msg)
		if err != nil {
			log.Error().Msg(err.Error())
			return
		}
	}
}
