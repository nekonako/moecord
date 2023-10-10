package websocket

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"
	"github.com/livekit/protocol/auth"
	"github.com/nekonako/moecord/config"
	"github.com/nekonako/moecord/infra"
	"github.com/nekonako/moecord/internal/sfu"
	"github.com/nekonako/moecord/pkg/api"
	"github.com/nekonako/moecord/pkg/middleware"
	"github.com/nekonako/moecord/pkg/tracer"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"
)

type wsConn struct {
	net.Conn
	// websocket connection hold user_id to identify connection
	UserID ulid.ULID `json:"user_id"`
}

type ConnMap struct {
	*sync.RWMutex

	// first key : channel_id
	// second key : user_id
	// hold user that associate to channel
	// this is used to broadcast message to all user that associate to specific channel
	channels map[string]map[string]*wsConn

	// this is map to detect list of online users that connected to websocket
	// key : user_id
	users map[string]bool

	// first key : server_id
	// second key : user_id
	// hold user that associate to server
	// this is used to broadcast message to all user that associate to spesific server
	servers map[string]map[string]*wsConn
}

var (
	connOnce sync.Once
	connMap  *ConnMap
)

type Websocket struct {
	config *config.Config
	infra  *infra.Infra
	sfu    *sfu.SFU
}

type channel struct {
	ID   ulid.ULID `db:"id"`
	Name string    `db:"name"`
}

type server struct {
	ID ulid.ULID `db:"id"`
}

type typingMessage struct {
	ChannelID ulid.ULID `json:"channel_id"`
	UserID    ulid.ULID `json:"user_id"`
	ServerID  ulid.ULID `json:"server_id"`
}

type userConnectionState struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

func New(c *config.Config, infra *infra.Infra, sfu *sfu.SFU) *Websocket {
	connOnce.Do(func() {
		connMap = &ConnMap{
			RWMutex:  &sync.RWMutex{},
			channels: make(map[string]map[string]*wsConn),
			users:    make(map[string]bool),
			servers:  make(map[string]map[string]*wsConn),
		}
	})
	return &Websocket{
		infra:  infra,
		config: c,
		sfu:    sfu,
	}
}

func (ws *Websocket) InitRouter(r *mux.Router) {

	v1 := r.NewRoute().Subrouter()
	v1.Use(middleware.Authentication(ws.config))
	v1.HandleFunc("/v1/channel/{channel_id}/voice/token", ws.creeteVoiceRoomToken).Methods(http.MethodGet)

	r.HandleFunc("/ws/voice/{channel_id}", ws.AcceptVoiceConnection)
	r.HandleFunc("/ws", ws.acceptConnection)

}

func (cat *Websocket) AcceptVoiceConnection(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "websocket.AcceptVoiceConnection")
	defer tracer.Finish(span)

	c, err := r.Cookie("access_token")
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return
	}

	accessToken := c.Value
	claim, err := middleware.ValidateToken(accessToken, cat.config.JWT.PrivateKey)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return
	}

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return
	}

	channelID := mux.Vars(r)["channel_id"]
	cid, _ := ulid.Parse(channelID)

	userIDStr := claim["sub"].(string)
	uid, _ := ulid.Parse(userIDStr)

	channel, err := cat.getChannelByID(ctx, uid, cid)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return
	}

	cat.sfu.CreateOrGetRoom(conn, channelID, channel.Name)

}

func (cat *Websocket) acceptConnection(w http.ResponseWriter, r *http.Request) {

	ctx, span := tracer.Start(r.Context(), "websocket.AcceptConnection")
	defer tracer.Finish(span)

	c, err := r.Cookie("access_token")
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return
	}

	accessToken := c.Value
	claim, err := middleware.ValidateToken(accessToken, cat.config.JWT.PrivateKey)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return
	}

	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return
	}

	userIDStr := claim["sub"].(string)
	userID, _ := ulid.Parse(userIDStr)
	channels, err := cat.getUserTextChannel(ctx, userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return
	}

	servers, err := cat.getUserServer(ctx, userID)
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
		return
	}

	connMap.Lock()
	con := &wsConn{
		UserID: userID,
		Conn:   conn,
	}
	chanArr := []string{}
	serverArr := []string{}
	for _, v := range channels {
		_, connChanExists := connMap.channels[v.ID.String()]
		if !connChanExists {
			connMap.channels[v.ID.String()] = make(map[string]*wsConn)
		}
		connMap.channels[v.ID.String()][userID.String()] = con
		chanArr = append(chanArr, v.ID.String())
	}

	for _, v := range servers {
		_, connServerExists := connMap.servers[v.ID.String()]
		if !connServerExists {
			connMap.servers[v.ID.String()] = make(map[string]*wsConn)
		}
		connMap.servers[v.ID.String()][userIDStr] = con
		newCOnnection := api.WebSocketMessage[any]{
			EventID: "NEW_CONNECTION",
			Data: userConnectionState{
				UserID: userIDStr,
				Status: "online",
			},
		}
		bm, _ := json.Marshal(newCOnnection)
		for _, userConn := range connMap.servers[v.ID.String()] {
			wsutil.WriteServerMessage(userConn.Conn, ws.OpText, bm)
		}
		serverArr = append(serverArr, v.ID.String())
	}
	connMap.Unlock()

	go cat.listenConnection(serverArr, chanArr, userIDStr, conn)

}

func (cat *Websocket) creeteVoiceRoomToken(w http.ResponseWriter, r *http.Request) {

	ctx, span := tracer.Start(r.Context(), "websocket.createVoiceRoomToken")
	defer tracer.Finish(span)

	at := auth.NewAccessToken(cat.config.LiveKit.ApiKey, cat.config.LiveKit.ApiSecret)
	channelID := mux.Vars(r)["channel_id"]
	userID := ctx.Value(middleware.Claim("user_id")).(string)
	grant := &auth.VideoGrant{
		RoomJoin:   true,
		RoomCreate: true,
		Room:       channelID,
	}
	at.AddGrant(grant).
		SetIdentity(userID).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		tracer.SpanError(span, err)
		log.Error().Ctx(ctx).Msg(err.Error())
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

func (w *Websocket) listenConnection(serverID, channelID []string, userID string, c net.Conn) {
	ctx, span := tracer.Start(context.Background(), "websocket.listenConnection")
	defer func() {
		tracer.Finish(span)
		c.Close()
	}()
	for {
		m, _, err := wsutil.ReadClientData(c)
		if err != nil {
			tracer.SpanError(span, err)
			log.Error().Ctx(ctx).Msg(err.Error())
			connMap.Lock()
			delete(connMap.users, userID)
			for _, v := range serverID {
				delete(connMap.servers[v], userID)
				for _, vv := range connMap.servers[v] {
					x := api.WebSocketMessage[any]{
						EventID: "USER_DISCONNECTED",
						Data: userConnectionState{
							UserID: userID,
							Status: "offline",
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
			tracer.SpanError(span, err)
			log.Error().Ctx(ctx).Msg(err.Error())
			continue
		}
		if wm.EventID == "STOP_TYPING" {
			message := api.WebSocketMessage[typingMessage]{}
			err = json.Unmarshal(m, &message)
			if err != nil {
				tracer.SpanError(span, err)
				log.Error().Ctx(ctx).Msg(err.Error())
				continue
			}
			for _, v := range connMap.servers[message.Data.ServerID.String()] {
				wsutil.WriteServerMessage(v.Conn, 0x1, m)
			}
		}

		if wm.EventID == "TYPING" {
			message := api.WebSocketMessage[typingMessage]{}
			err = json.Unmarshal(m, &message)
			if err != nil {
				tracer.SpanError(span, err)
				log.Error().Ctx(ctx).Msg(err.Error())
				continue
			}
			for _, v := range connMap.servers[message.Data.ServerID.String()] {
				wsutil.WriteServerMessage(v.Conn, 0x1, m)
			}
		}
	}
}

func (ws *Websocket) GetOnlineUserByServerID(ctx context.Context, userID string) {

}
