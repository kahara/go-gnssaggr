package reader

import (
	"encoding/json"
	"time"
)

func parser(lines <-chan []byte, reports chan<- Report) {
	var (
		err  error
		tpv  TPV
		sky  SKY
		line []byte
	)

	for {
		select {
		case line = <-lines:
			// Try the more commong TPV report first
			if err = json.Unmarshal(line, &tpv); err != nil || tpv.Class != "TPV" {
				// Try a SKY report
				if err = json.Unmarshal(line, &sky); err == nil && sky.Class == "SKY" {
					// This is a valid SKY report
					sky.Time = time.Now()
					reports <- sky
				}
			} else {
				// This is a valid TPV report
				reports <- tpv
			}
		}
	}
}
