package utils

import (
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type getRet string

func (ret getRet) String() (string, error) {
	return string(ret), nil
}

func (ret getRet) Int() (int, error) {
	return strconv.Atoi(string(ret))
}

func (ret getRet) Int64() (int64, error) {
	return strconv.ParseInt(string(ret), 10, 64)
}

func (ret getRet) Uint() (uint, error) {
	newRet, err := strconv.ParseUint(string(ret), 10, 32)
	return uint(newRet), err
}

func (ret getRet) Uint64() (uint64, error) {
	return strconv.ParseUint(string(ret), 10, 64)
}

func (ret getRet) Bool() (bool, error) {
	return strconv.ParseBool(string(ret))
}

func GetQuery(r *http.Request, key string) getRet {
	return getRet(r.URL.Query().Get(key))
}

func GetQueryDefault(r *http.Request, key string, defaultValue string) getRet {
	value := r.URL.Query().Get(key)
	if value == "" {
		value = defaultValue
	}
	return getRet(value)
}

func GetParam(r *http.Request, key string) getRet {
	return getRet(mux.Vars(r)[key])
}

func ClientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}

	if addr := r.Header.Get("X-Appengine-Remote-Addr"); addr != "" {
		return addr
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
