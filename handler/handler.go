package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/wpjunior/hyperbigbang/fetch"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url, err := convertPathToURL(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accept := r.Header.Get("accept")
	if accept == "" {
		accept = "application/json"
	} else {
		r.Header.Del("accept")
	}

	includes := extractRelIncludes(r.URL)
	result, err := fetch.URL(url, r.Header, includes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if accept == "application/msgpack" {
		writeMsgPackResponse(w, result)
	} else {
		writeJSONResponse(w, result)
	}
}

func writeJSONResponse(w http.ResponseWriter, result interface{}) {
	w.Header().Add("content-type", "application/json")
	err := json.NewEncoder(w).Encode(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func writeMsgPackResponse(w http.ResponseWriter, result interface{}) {
	w.Header().Add("content-type", "application/msgpack")
	err := msgpack.NewEncoder(w).Encode(result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func extractRelIncludes(u *url.URL) []string {
	raw := u.Query().Get("include")
	return strings.Split(raw, ",")
}
