package websocket

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/middleware"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type conn struct {
	*sync.RWMutex
	conn map[string][]net.Conn
}

var (
	connOnce sync.Once
	connMap  *conn
)

type websocket struct {
	config *config.Config
	infra  *infra.Infra
}

type channel struct {
	ID ulid.ULID `db:"channel_id"`
}

type Message struct {
	ID        ulid.ULID    `json:"id"`
	ChannelID ulid.ULID    `json:"channel_id"`
	SenderID  ulid.ULID    `json:"sender_id"`
	Content   string       `json:"content"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

func New(c *config.Config, infra *infra.Infra) *websocket {

	connOnce.Do(func() {
		connMap = &conn{
			conn:    make(map[string][]net.Conn),
			RWMutex: &sync.RWMutex{},
		}
	})

	return &websocket{
		infra:  infra,
		config: c,
	}

}

func (ws *websocket) InitRouter(r *mux.Router) {
	r.HandleFunc("/ws", ws.AcceptConnection)
	go ws.Sub()
}

func (w *websocket) GetUserChannel(ctx context.Context, userID ulid.ULID) ([]channel, error) {
	ctx, span := tracer.Start(ctx, "websocket.GetUserChannel")
	defer tracer.Finish(span)

	result := []channel{}
	query := `SELECT channel_id FROM channel_member WHERE user_id = $1`
	err := w.infra.Postgres.SelectContext(ctx, &result, query, userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Msg(err.Error())
		return result, err
	}

	return result, nil

}

func (cat *websocket) AcceptConnection(w http.ResponseWriter, r *http.Request) {

	ctx, span := tracer.Start(r.Context(), "websocket.AcceptConnection")
	defer tracer.Finish(span)

	c, err := r.Cookie("access_token")
	if err != nil {
		api.NewHttpResponse().
			WithCode(http.StatusBadGateway).
			WithError(err)
		return
	}

	accessToken := c.Value
	claim, err := middleware.ValidateToken(accessToken, cat.config.JWT.PrivateKey)
	if err != nil {
		api.NewHttpResponse().
			WithCode(http.StatusBadGateway).
			WithError(err)
		return
	}

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		api.NewHttpResponse().
			WithCode(http.StatusBadGateway).
			WithError(err)
		return
	}

	userIDStr := claim["sub"].(string)
	userID, _ := ulid.Parse(userIDStr)
	channels, err := cat.GetUserChannel(ctx, userID)
	if err != nil {
		api.NewHttpResponse().
			WithCode(http.StatusBadGateway).
			WithError(err)
		return
	}

	connMap.Lock()
	for _, v := range channels {
		connMap.conn[v.ID.String()] = append(connMap.conn[v.ID.String()], conn)
	}
	connMap.Unlock()

}

func (w *websocket) Sub() {
	w.infra.Nats.Subscribe("NEW_CHANNEL_MESSAGE", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
		message := Message{}

		err := json.Unmarshal(m.Data, &message)
		if err != nil {
			log.Error().Msg(err.Error())
		}

		if c, ok := connMap.conn[message.ChannelID.String()]; ok {
			go w.broadcastChannel(m.Data, c)
		}
	})
}

func (w *websocket) broadcastChannel(m []byte, conn []net.Conn) {
	for _, c := range conn {
		wsutil.WriteServerMessage(c, 0x1, m)
	}
}
