package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/token"
	"github.com/gost-codes/sweet_dreams/util"
)

type Server struct {
	store      db.Store
	config     util.Config
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	server := &Server{store: store, config: config}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	server.router = router
}

func (server *Server) Start(addr *string) {
	if addr == nil {

		fmt.Println("Server started\n Listening on:", 8080)
	} else {

		fmt.Println("Server started\n Listening on:", addr)
	}
	server.router.Run()
}
