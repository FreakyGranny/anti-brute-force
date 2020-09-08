package storage

import (
	"fmt"
	"strconv"
)

func sslModeToString(sslEnable bool) string {
	if sslEnable {
		return "enable"
	}

	return "disable"
}

// BuildDsn build dsn string from params.
func BuildDsn(host string, port int, user string, password string, dbName string, sslMode bool) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host,
		strconv.Itoa(port),
		user,
		dbName,
		password,
		sslModeToString(sslMode),
	)
}
