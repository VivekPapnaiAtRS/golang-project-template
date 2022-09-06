package server

import (
	"VivekPapnaiAtRS/template/utils"
	"net/http"
)

func (srv *Server) greet(resp http.ResponseWriter, req *http.Request) {
	utils.EncodeJSONBody(resp, http.StatusOK, map[string]interface{}{
		"message": "Hello Bro",
	})
}
