package logger

import (
	"log"
	"os"
	"strings"
	"sync"
)

var (
	debugEnabled bool
	once         sync.Once
)

func DebugEnabled() bool {
	once.Do(func() {
		v := strings.ToLower(strings.TrimSpace(os.Getenv("DEBUG")))
		debugEnabled = v == "1" || v == "true" || v == "yes" || v == "on"
	})
	return debugEnabled
}

func Debugf(format string, args ...any) {
	if !DebugEnabled() {
		return
	}
	log.Printf("DEBUG "+format, args...)
}

func Infof(format string, args ...any) {
	log.Printf("INFO "+format, args...)
}

func Errorf(format string, args ...any) {
	log.Printf("ERROR "+format, args...)
}
