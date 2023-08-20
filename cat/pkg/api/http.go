package api

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type HttpResponse struct {
	Code     int         `json:"code"`
	Message  string      `json:"message,omitempty"`
	Error    string      `json:"error,omitempty"`
	Errors   any         `json:"errors,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	MetaData *MetaData   `json:"metadata,omitempty"`
}

type MetaData struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{}
}

func (r *HttpResponse) WithCode(code int) *HttpResponse {
	r.Code = code
	return r
}

func (r *HttpResponse) WitMessage(m string) *HttpResponse {
	r.Message = m
	return r
}

func (r *HttpResponse) WithError(e error) *HttpResponse {
	r.Error = e.Error()
	return r
}

func (r *HttpResponse) WithErrors(e any) *HttpResponse {
	r.Errors = e
	return r
}

func (r *HttpResponse) WithData(data interface{}) *HttpResponse {
	r.Data = data
	return r
}

func (r *HttpResponse) WithMetaData(m *MetaData) *HttpResponse {
	m.TotalPage = (m.Total + m.PerPage - 1) / m.PerPage
	r.MetaData = m
	return r
}

func (r *HttpResponse) SendJSON(w http.ResponseWriter) {

	b, err := json.Marshal(r)
	if err != nil {
		log.Error().Err(err).Msg("failed encode response")
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	w.Write(b)
}
