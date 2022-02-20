package reader

import (
	"bufio"
	"fmt"
	"github.com/cespare/xxhash"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

const (
	backoffDefault       = 1.0
	backoffMultiplier    = 1.35
	backoffLimit         = 20
	backoffWhenConnected = time.Duration(5 * time.Second)             // Sleep fixed amount of time when errors happen after connecting
	trigger              = "?WATCH={\"enable\":true,\"json\":true}\n" // https://gpsd.gitlab.io/gpsd/gpsd_json.html#_watch
)

func Read(config Config, reports chan<- Report) {

	var (
		err      error
		hostport = fmt.Sprintf("%s:%d", config.Host, config.Port)
		backoff  = backoffDefault
		conn     net.Conn
		n        int
		scanner  *bufio.Scanner
		lines    = make(chan []byte, 2)
		lineHash uint64
		prevHash uint64
		line     []byte
	)

	// Newline-delimited chunks are fed to parser
	go parser(lines, reports)

	for {
		if conn, err = net.Dial("tcp", hostport); err != nil {
			log.Error().Err(err).Msgf("Error while dialing %s, sleeping for %fs before starting over", hostport, backoff)
			time.Sleep(time.Duration(backoff) * time.Second)

			// Sleep longer after next error
			backoff = backoff * backoffMultiplier
			if backoff > backoffLimit {
				backoff = backoffLimit
			}
			continue
		}

		log.Printf("Connected to %s", hostport)

		// Reset backoff to default
		backoff = backoffDefault

		// Ask gpsd to start producing JSON reports
		if n, err = fmt.Fprint(conn, trigger); err != nil {
			log.Error().Err(err).Msgf("Error while asking gpsd to start producing JSON reports on %s, %d/%d bytes written; sleeping for %fs before starting over", hostport, n, len(trigger), backoffWhenConnected.Seconds())
			time.Sleep(backoffWhenConnected)
			break
		}

		// Consume newline-separated chunks of bytes from the TCP stream
		scanner = bufio.NewScanner(conn)
		for {
			if scanner.Scan() {
				line = scanner.Bytes()

				// Skip consecutive duplicates
				lineHash = xxhash.Sum64(line)
				if lineHash != prevHash {
					lines <- line
					prevHash = lineHash
				}
			} else {
				log.Error().Err(scanner.Err()).Msgf("Error while scanning lines from %s, sleeping for %fs before starting over", hostport, backoffWhenConnected.Seconds())
				time.Sleep(backoffWhenConnected)
				break
			}
		}
	}
}
