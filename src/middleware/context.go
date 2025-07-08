package middleware

import (
	"RestGoTest/src/util"
	"context"
	"fmt"
	"net/http"
	"time"
)

func ContextDelayAbortMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for i := 0; i < 5; i++ {
			select {
			case <-r.Context().Done():
				fmt.Println("DELETE Services Cancel By Client")
				util.ResponseWithError(w, http.StatusRequestTimeout, "درخواست کنسل شد")
				return
			case <-time.After(1 * time.Second):
			}
		}
		next.ServeHTTP(w, r)
	})
}

func ContextAbortMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			fmt.Println("Client Cancel Call-Service")
			http.Error(w, "درخواست کنسل شد", http.StatusRequestTimeout)
			return
		default:
			next.ServeHTTP(w, r)
		}
	})
}

func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
