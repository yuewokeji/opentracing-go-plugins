package othttp

import (
	"net"
	"net/http"
	"strings"
)

func GetRemoteAddr(r *http.Request) string {
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
