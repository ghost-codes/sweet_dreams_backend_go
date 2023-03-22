package api

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/gost-codes/sweet_dreams/db/sqlc"
	"github.com/gost-codes/sweet_dreams/token"
	"github.com/gost-codes/sweet_dreams/util"
	"github.com/gost-codes/sweet_dreams/worker"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if binder, ok := binding.Validator.Engine().(*validator.Validate); ok {
		binder.RegisterValidation("bookingType", validateBookingType)
	}

	//User router
	router.POST("/users/sign_up", server.createUserWithEmailPassword)
	router.POST("/users/login", server.loginWithEmailPassword)
	router.POST("/users/socials/sigin", server.signInUserSocial)
	router.GET("/verify_email", server.verifyEmail)

	//Admin router
	router.POST("/admin/login", server.adminLogin)

	//with auth middleware
	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker, &server.store))
	authRouter.POST("/send_verification_email", server.sendVerificationEmail)
	authRouter.GET("/nurses", server.fetchNurses)
	authRouter.POST("/bookings/request", server.createBooking)
	authRouter.POST("/bookings/approvals", server.userApprovals)
	authRouter.GET("/bookings/request", server.listUserBookingReqs)
	authRouter.GET("/bookings/request/:id", server.bookingReq)
	authRouter.DELETE("/bookings/request/:id", server.deleteBookingReq)

	superAdminRouter := router.Group("/admin").Use(adminAuthMiddleware(server.tokenMaker, &server.store, true))
	superAdminRouter.POST("/create", server.createAdmin)

	adminRouter := router.Group("/admin").Use(adminAuthMiddleware(server.tokenMaker, &server.store, false))
	adminRouter.POST("/nurses/create", server.createNurse)
	adminRouter.GET("/nurses/:id", server.fetchNurse)
	adminRouter.GET("/nurses", server.fetchNurses)
	adminRouter.DELETE("/nurses/:id", server.deleteNurse)
	adminRouter.GET("/bookings/request", server.adminBookings)
	adminRouter.GET("/bookings/request/:id", server.adminBookingsByID)
	adminRouter.POST("/bookings/approvals", server.adminApproveBooking)
	adminRouter.GET("/bookings/approvals", server.adminGetApprovals)
	adminRouter.DELETE("/bookings/approvals/:id", server.deleteApprovals)
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

func genericResponse(msg string) gin.H {
	return gin.H{"error": msg}
}
