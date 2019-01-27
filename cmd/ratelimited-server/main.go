package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	rls "github.com/kazunobu-fujii/ratelimited-server"
	"golang.org/x/time/rate"
)

const (
	defaultTokens    = 1
	defaultCapSecond = 1
	defaultPort      = "localhost:8080"
	defaultChunkSize = 8
	defaultWait      = 10
	defaultTimeout   = 30
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, "PANIC: %v\n", err)
			os.Exit(1)
		}
	}()

	var requestcap int
	flag.IntVar(&requestcap, "r", defaultCapSecond, "request rate limit (per second)")
	var tokens int
	flag.IntVar(&tokens, "d", defaultTokens, "bucket depth")

	var timeout int
	s := rls.NewServer()
	flag.StringVar(&s.Config.Path, "p", "response.dat", "response body file path")
	flag.StringVar(&s.Config.Ctype, "t", "text/html; charset=UTF-8", "content type")
	flag.IntVar(&s.Config.Chunk, "c", defaultChunkSize, "chunk size")
	flag.IntVar(&s.Config.Wait, "w", defaultWait, "chunk delay (ms)")
	flag.IntVar(&timeout, "o", defaultTimeout, "timeout (s)")
	flag.StringVar(&s.Config.Port, "s", defaultPort, "listening server")
	flag.Parse()

	s.Config.Timeout = time.Duration(timeout) * time.Second
	s.Config.RateLimiter = rls.NewMultiLimiter(
		rate.NewLimiter(per(requestcap, time.Second), tokens),
	)

	if err := s.Main(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func per(eventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(eventCount))
}
