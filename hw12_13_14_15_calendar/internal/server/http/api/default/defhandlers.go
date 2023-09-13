package api

import (
	"net/http"

	iface "github.com/exiffM/otus_homework/hw12_13_14_15_calendar/internal/interface"
	"github.com/mailru/easyjson"
)

type DefaultHandler struct {
	Response DefaultResponse
	Logger   iface.Logger
}

func (ah *DefaultHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = r
	_, _, err := easyjson.MarshalToHTTPResponseWriter(ah.Response, w)
	if err != nil {
		ah.Logger.Info("Write http response error on default handler")
		return
	}
}
