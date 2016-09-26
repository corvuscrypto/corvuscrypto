package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var secretKey []byte

func checkAuth(f httprouter.Handle) httprouter.Handle {
	return func(_w http.ResponseWriter, _r *http.Request, _p httprouter.Params) {
		//code to check auth
		var authorized bool
		if authorized {
			f(_w, _r, _p)
		} else {
			_w.WriteHeader(401)
		}
	}
}
