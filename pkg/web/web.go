package web

import (
	"context"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ejuju/go-utils/pkg/logs"
)

// Access logging middleware logs incoming HTTP requests
func AccessLoggingMiddleware(logger logs.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resrec := &httpStatusRecorder{ResponseWriter: w} // use custom response writer to record status
			before := time.Now()                             // record timestamp before request is handled
			h.ServeHTTP(resrec, r)                           //
			dur := time.Since(before)                        // calculate duration to handle request

			// Log
			logstr := fmt.Sprintf("%d %-4s %5dμs %s", resrec.statusCode, r.Method, dur.Microseconds(), r.URL.Path)
			logger.Log(logstr)
		})
	}
}

type httpStatusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (srec *httpStatusRecorder) WriteHeader(statusCode int) {
	srec.statusCode = statusCode
	srec.ResponseWriter.WriteHeader(statusCode)
}

type PanicHandler func(w http.ResponseWriter, r *http.Request, err any)

// Panic recovery middleware logs the recovered error and executes the onPanic callback function.
func PanicRecoveryMiddleware(onPanic PanicHandler) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if onPanic != nil {
						onPanic(w, r, err)
					}
				}
			}()

			h.ServeHTTP(w, r)
		})
	}
}

func NewServerWithDefaults(h http.Handler, port int) *http.Server {
	out := &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		Handler:           h,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		MaxHeaderBytes:    8000,
	}
	return out
}

// Listens for incoming connections and
// await interrupt signal (or server error) for graceful shutdown
func RunServer(s *http.Server) error {
	// Start HTTP server in seperate goroutine
	errc := make(chan error, 1)
	go func() { errc <- s.ListenAndServe() }()

	// Await interrupt signal (or server error) for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-done:
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		shutdownErr := s.Shutdown(ctx)
		if shutdownErr != nil {
			panic(shutdownErr)
		}
	case err := <-errc:
		return err
	}
	return nil
}

func PermanentRedirectHandler(toURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, toURL, http.StatusPermanentRedirect) }
}

func VisitorHash(r *http.Request) string {
	// get true IP
	ipAddr := r.Header.Get("X-Forwarded-For")
	if ipAddr == "" {
		ipAddr = net.ParseIP(r.RemoteAddr).String()
	}
	// Hash IP addr and user-agent
	hash := sha1.New()
	_, err := hash.Write([]byte(ipAddr + r.UserAgent()))
	if err != nil {
		panic(err)
	}
	// Return base64 hex encoded hash
	return base32.StdEncoding.EncodeToString(hash.Sum(nil))
}

func ReadAndServeFile(path string) http.HandlerFunc {
	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", http.DetectContentType(raw))
		w.WriteHeader(http.StatusOK)
		w.Write(raw)
	}
}

func ServeRaw(v []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if w.Header().Get("Content-Type") == "" {
			w.Header().Set("Content-Type", http.DetectContentType(v))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(v)
	}
}

func FileServer(dirPath, httpPrefix string) http.Handler {
	return http.StripPrefix(httpPrefix, http.FileServer(http.Dir(dirPath)))
}
