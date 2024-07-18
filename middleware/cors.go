package middleware

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

type CorsMiddleware struct {
}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}

func (m *CorsMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Add("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			httpx.Ok(w)

			return
		}

		next(w, r)
	}
}
