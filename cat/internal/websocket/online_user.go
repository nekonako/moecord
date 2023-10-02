package websocket

import "github.com/oklog/ulid/v2"

func GetOnlineUser(ws *Websocket, serverID ulid.ULID) map[string]bool {
	res := map[string]bool{}
	connMap.Lock()
	for _, v := range connMap.servers[serverID.String()] {
		res[v.UserID.String()] = true
	}
	connMap.Unlock()
	return res
}
