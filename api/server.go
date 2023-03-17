package api

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/token"
	"github.com/gost-codes/sweet_dreams/util"
	"github.com/gost-codes/sweet_dreams/worker"
)

type Server struct {
	store           db.Store
	config          util.Config
	tokenMaker      token.Maker
	router          *gin.Engine
	firebase        *firebase.App
	taskDistributor worker.TaskDistributor
}

func NewServer(store db.Store, config util.Config, distributor worker.TaskDistributor) (*Server, error) {
	server := &Server{store: store, config: config}

	tokenMaker, err := token.NewPasetoMaker(config.SecretKey)

	if err != nil {
		return nil, fmt.Errorf("Cannot create token maker: %v", err)
	}

	server.tokenMaker = tokenMaker
	app, err := util.InitializeFirebaseApp(context.Background(), nil)
	server.taskDistributor = distributor

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
	router.POST("/users/socials/sigin", server.signInUserSocial)
	router.GET("/verify_email", server.verifyEmail)

	//with auth middleware
	authRouter := router.Use(authMiddleware(server.tokenMaker))

	authRouter.POST("/send_verification_email", server.sendVerificationEmail)
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
