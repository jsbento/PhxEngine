package logger

import (
	"log"

	"github.com/davecgh/go-spew/spew"
)

func SpewLog(v interface{}, customMsg ...string) {
	msg := ""
	if len(customMsg) > 0 {
		msg = customMsg[0]
	}
	log.Printf("%s: %s", msg, spew.Sdump(v))
}
