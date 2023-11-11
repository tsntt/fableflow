package middleware

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/tsntt/fableflow/src/util"
	"github.com/uptrace/bunrouter"
	"golang.org/x/time/rate"
)

func RateLimit(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 10*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return func(w http.ResponseWriter, r bunrouter.Request) error {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			util.WriteJson(w, http.StatusInternalServerError, "Something went wrong")
			return errors.New("could not parse addr")
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			util.WriteJson(w, http.StatusTooManyRequests, "to many requests")
			return nil
		}

		mu.Unlock()

		return next(w, r)
	}
}
