package util

import (
	"log"
)

var logEnabled bool

func EnableLog() { logEnabled = true }

func Log(s ...interface{}) {
	if !logEnabled {
		return
	}

	v := []interface{}{"[ZORM]"}
	v = append(v, s...)

	log.Println(v...)
}
