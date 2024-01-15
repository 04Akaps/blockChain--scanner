package log

import (
	"log"
	"os"
	"strings"
)

var logFile *os.File

func SetLog(p string) {

	if !strings.HasSuffix(p, ".txt") {
		panic(".txt is not suffixed at logName Env")
	} else {
		if f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
			if os.IsNotExist(err) {
				if f, err = os.Create(p); err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		} else {
			logFile = f
			log.SetOutput(f)
		}
	}

}

func GetLogFile() *os.File {
	return logFile
}

func InfoLog(w string) {
	log.Printf("[INFO] " + w)
}

func ErrLog(w string) {
	log.Printf("[ERR] " + w)
}

func CritLog(w string) {
	log.Printf("[Crit] " + w)
	panic(w)
}
