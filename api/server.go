package api

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
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
	firebase   *firebase.App
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	server := &Server{store: store, config: config}

	tokenMaker, err := token.NewPasetoMaker(config.SecretKey)

	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %v", err)
	}

	server.tokenMaker = tokenMaker
	app, err := util.InitializeFirebaseApp(context.Background(), nil)

	if err != nil {
		return nil, err
	}
	server.firebase = app
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users/sign_up", server.createUserWithEmailPassword)
	router.POST("/users/login", server.loginWithEmailPassword)
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

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
