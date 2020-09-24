package storage

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

func sslModeToString(sslEnable bool) string {
	if sslEnable {
		return "enable"
	}

	return "disable"
}

// BuildDsn build dsn string from params.
func BuildDsn(host string, port int, user string, password string, dbName string, sslEnable bool) string {
	u := &url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(user, password),
		Host:     net.JoinHostPort(host, strconv.Itoa(port)),
		Path:     dbName,
		RawQuery: strings.Join([]string{"sslmode", sslModeToString(sslEnable)}, "="),
	}

	return u.String()
}
