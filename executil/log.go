package executil

import "go.gary.cool/go-stockutil/log"

func LogOutput(outlevel log.Level, errlevel log.Level) OutputLineFunc {
	return func(line string, err bool) {
		if err {
			log.Log(outlevel, line)
		} else {
			log.Log(errlevel, line)
		}
	}
}

func LogOutputFunc(line string, err bool) {
	if err {
		log.Error(line)
	} else {
		log.Debug(line)
	}
}
