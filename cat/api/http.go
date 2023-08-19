package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Response struct {
	Code     int         `json:"code"`
	Message  string      `json:"message"`
	Error    string      `json:"error,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	MetaData *MetaData   `json:"metadata,omitempty"`
}

type MetaData struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) WithCode(code int) *Response {
	r.Code = code
	return r
}

func (r *Response) WitMessage(m string) *Response {
	r.Message = m
	return r
}

func (r *Response) WithError(e string) *Response {
	r.Error = e
	return r
}

func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) WithMetaData(m *MetaData) *Response {
	m.TotalPage = (m.Total + m.PerPage - 1) / m.PerPage
	r.MetaData = m
	return r
}

func (r *Response) SendJSON(w http.ResponseWriter) {

	b, err := json.Marshal(r)
	if err != nil {
		log.Error().Err(err).Msg("failed encode response")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	w.Write(b)
}
