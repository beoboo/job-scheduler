package server

import (
	"fmt"
	"net/http"
)

type HttpBasicAuthService struct {
	next http.Handler
}

func NewHttpBasicService(next http.Handler) *HttpBasicAuthService {
	return &HttpBasicAuthService{
		next: next,
	}
}

func (h *HttpBasicAuthService) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Checking auth: %s\n", req.Body)

	user, pass, ok := req.BasicAuth()

	if !ok || !checkCredentials(user, pass) {
		rw.WriteHeader(401)
		rw.Write([]byte("Unauthorised.\n"))
		return
	}

	h.next.ServeHTTP(rw, req)
}

func checkCredentials(user, pwd string) bool {
	return user == "user" && pwd == "test"
}
