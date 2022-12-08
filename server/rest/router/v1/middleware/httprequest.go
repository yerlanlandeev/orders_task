package middleware

import (
	"net/http"
	"strings"
)

type HttpMethod struct {
	method string
}

func NewHttpMethod(method string) *HttpMethod {
	return &HttpMethod{
		method: method,
	}
}

func (h HttpMethod) OnlyThisMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h.method != r.Method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func OnlyContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok := strings.EqualFold(r.Header.Get("Content-Type"), "application/json"); !ok {
			http.Error(w, "content type must be", http.StatusPreconditionFailed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
