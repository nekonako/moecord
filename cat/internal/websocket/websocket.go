package websocket

import (
	"context"
	"encoding/json"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"
	"github.com/livekit/protocol/auth"
	"github.com/nats-io/nats.go"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/internal/message/usecase"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/middleware"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"sync"
	"time"
)

type conn struct {
	*sync.RWMutex
	channels map[string]map[string]*myConn
	users    map[string]bool
	servers  map[string]map[string]*myConn
}

type myConn struct {
	net.Conn
	userid ulid.ULID
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
	ID ulid.ULID `db:"id"`
}

type server struct {
	ID ulid.ULID `db:"id"`
}

type TypingMessage struct {
	ChannelID ulid.ULID `json:"channel_id"`
	UserID    ulid.ULID `json:"user_id"`
	ServerID  ulid.ULID `json:"server_id"`
}

func New(c *config.Config, infra *infra.Infra) *websocket {
	connOnce.Do(func() {
		connMap = &conn{
			RWMutex:  &sync.RWMutex{},
			channels: make(map[string]map[string]*myConn),
			users:    make(map[string]bool),
			servers:  make(map[string]map[string]*myConn),
		}
	})
	return &websocket{
		infra:  infra,
		config: c,
	}
}

func (ws *websocket) InitRouter(r *mux.Router) {
	r.HandleFunc("/v1/room/token", ws.CreeteVoiceRoomToken).Methods(http.MethodGet)
	r.HandleFunc("/ws", ws.AcceptConnection)
	go ws.Sub()
}

func (w *websocket) GetUserChannel(ctx context.Context, userID ulid.ULID) ([]channel, error) {
	ctx, span := tracer.Start(ctx, "websocket.GetUserChannel")
	defer tracer.Finish(span)

	result := []channel{}
	query := `
       SELECT
         distinct(c.id)
       FROM channel AS c
       LEFT JOIN channel_member AS cm ON cm.channel_id = c.id
       WHERE (c.is_private = false OR cm.user_id = $1)
    `
	err := w.infra.Postgres.SelectContext(ctx, &result, query, userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return result, err
	}

	return result, nil

}

func (w *websocket) GetUserServer(ctx context.Context, userID ulid.ULID) ([]server, error) {
	ctx, span := tracer.Start(ctx, "websocket.GetUserChannel")
	defer tracer.Finish(span)

	result := []server{}
	query := `
       SELECT
        server_id as id
       FROM server_member WHERE user_id = $1
       `
	err := w.infra.Postgres.SelectContext(ctx, &result, query, userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
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

	servers, err := cat.GetUserServer(ctx, userID)
	if err != nil {
		api.NewHttpResponse().
			WithCode(http.StatusBadGateway).
			WithError(err)
		return
	}

	connMap.Lock()
	con := &myConn{
		userid: userID,
		Conn:   conn,
	}
	chanArr := []string{}
	serverArr := []string{}
	for _, v := range channels {
		_, connChanExists := connMap.channels[v.ID.String()]
		if !connChanExists {
			connMap.channels[v.ID.String()] = make(map[string]*myConn)
		}
		connMap.channels[v.ID.String()][userID.String()] = con
		chanArr = append(chanArr, v.ID.String())
	}

	for _, v := range servers {
		_, connServerExists := connMap.servers[v.ID.String()]
		if !connServerExists {
			connMap.servers[v.ID.String()] = make(map[string]*myConn)
		}
		connMap.servers[v.ID.String()][userIDStr] = con
		serverArr = append(serverArr, v.ID.String())
	}
	connMap.Unlock()

	go cat.ListenConnection(serverArr, chanArr, userIDStr, conn)

}

func (w *websocket) Sub() {
	w.infra.Nats.Subscribe("NEW_CHANNEL_MESSAGE", func(m *nats.Msg) {
		message := api.WebSockerMessage[usecase.SaveMessageResponse]{}
		err := json.Unmarshal(m.Data, &message)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		if c, ok := connMap.channels[message.Data.ChannelID.String()]; ok {
			go w.broadcastChannel(message.Data, m.Data, c)
			m.AckSync()
		}
	})
}

func (w *websocket) broadcastChannel(m usecase.SaveMessageResponse, mb []byte, conn map[string]*myConn) {
	for _, c := range conn {
		if c.userid == m.SenderID {
			continue
		}
		wsutil.WriteServerMessage(c.Conn, 0x1, mb)
	}
}

func (cat *websocket) CreeteVoiceRoomToken(w http.ResponseWriter, r *http.Request) {

	at := auth.NewAccessToken(cat.config.LiveKit.ApiKey, cat.config.LiveKit.ApiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     "1",
	}
	at.AddGrant(grant).
		SetIdentity("xxxx").
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		log.Error().Msg(err.Error())
		api.NewHttpResponse().
			WithCode(http.StatusInternalServerError).
			WitMessage("internal server error").
			SendJSON(w)
		return
	}

	api.NewHttpResponse().WithData(map[string]string{
		"token": token,
	}).WithCode(200).SendJSON(w)

}

func (w *websocket) ListenConnection(serverID, channelID []string, userID string, c net.Conn) {
	defer func() {
		c.Close()
	}()
	for {
		m, _, err := wsutil.ReadClientData(c)
		if err != nil {
			log.Info().Msg("failed read data")
			connMap.Lock()
			delete(connMap.users, userID)
			for _, v := range serverID {
				delete(connMap.servers[v], userID)
				for _, vv := range connMap.servers[v] {
					x := api.WebSockerMessage[any]{
						EventID: "PAIR_CLOSED",
						Data: map[string]string{
							"user_id": userID,
						},
					}
					xx, _ := json.Marshal(x)
					wsutil.WriteServerMessage(vv.Conn, 0x1, xx)
				}
			}
			for _, v := range channelID {
				delete(connMap.channels[v], userID)
			}
			connMap.Unlock()
			return
		}
		wm := api.WebSocketMessage[any]{}
		err = json.Unmarshal(m, &wm)
		if err != nil {
			log.Error().Msg(err.Error())
			continue
		}
		if wm.EventID == "STOP_TYPING" {
			message := api.WebSocketMessage[TypingMessage]{}
			err = json.Unmarshal(m, &message)
			if err != nil {
				log.Error().Msg(err.Error())
				continue
			}
			for _, v := range connMap.servers[message.Data.ServerID.String()] {
				wsutil.WriteServerMessage(v.Conn, 0x1, m)
			}
		}

		if wm.EventID == "TYPING" {
			message := api.WebSocketMessage[TypingMessage]{}
			err = json.Unmarshal(m, &message)
			if err != nil {
				log.Error().Msg(err.Error())
				continue
			}
			for _, v := range connMap.servers[message.Data.ServerID.String()] {
				wsutil.WriteServerMessage(v.Conn, 0x1, m)
			}
		}
	}
}
