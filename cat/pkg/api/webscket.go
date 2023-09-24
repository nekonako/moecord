package api

type WebSockerMessage[T any] struct {
	EventID string `json:"event_id"`
	Data    T      `json:"data"`
}
