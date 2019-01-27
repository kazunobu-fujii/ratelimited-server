package rls

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Config ...
type Config struct {
	Path        string
	Ctype       string
	Chunk       int
	Wait        int
	Timeout     time.Duration
	Port        string
	RateLimiter RateLimiter
}

// Server ...
type Server struct {
	Logger *log.Logger
	Config Config
}

const logFlag = log.Ldate | log.Ltime | log.Lshortfile

// NewServer ...
func NewServer() *Server {
	return &Server{
		Logger: log.New(os.Stderr, "ratelimited-server ", logFlag),
	}
}

// Main ...
func (s *Server) Main() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Printf("request %v %v\n", r.Method, r.URL.Path)

		w.Header().Set("Content-Type", s.Config.Ctype)

		if err := s.service(context.Background(), w, r); err != nil {
			s.Logger.Printf("error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 Internal Server Error")
		}
	})

	s.Logger.Printf("Listening on %s", s.Config.Port)
	s.Logger.Fatal(http.ListenAndServe(s.Config.Port, nil))

	return nil
}

func (s *Server) service(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, s.Config.Timeout)
	defer cancel()

	if err := s.Config.RateLimiter.Wait(ctx); err != nil {
		return err
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("expected http.ResponseWriter to be http.Flusher")
	}

	file, err := os.Open(s.Config.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, s.Config.Chunk)
	for {
		n, err := file.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, string(buf[:n]))
		flusher.Flush()
		time.Sleep(time.Duration(s.Config.Wait) * time.Millisecond)
	}

	s.Logger.Printf("served %v %v\n", r.Method, r.URL.Path)
	return nil
}
