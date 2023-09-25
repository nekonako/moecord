package api

type WebSocketMessage[T any] struct {
	EventID string `json:"event_id"`
	Data    T      `json:"data"`
}
