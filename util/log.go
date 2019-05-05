package util

import (
	"log"
	"os"
	"strings"
)

var shouldLog = strings.ToLower(os.Getenv("ZORM_DEBUG_LOG")) == "true"

func Log(s ...interface{}) {
	if !shouldLog {
		return
	}

	v := []interface{}{"[ZORM]"}
	v = append(v, s...)

	log.Println(v...)
}
