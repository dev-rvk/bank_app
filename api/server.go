package api

import (
	"fmt"

	db "github.com/devrvk/simplebank/db/sqlc"
	"github.com/devrvk/simplebank/token"
	"github.com/devrvk/simplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// server struct contains a server router and store
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// constructor for the server struct with routes initialized
func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker error: %w", err)
	}

	// custom validator for currency
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// takes store from arguments and uses gin Default router
	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	// returns server object of type Server contains store and router
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	// define route
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.POST("/transfers", server.createTransfer)

	// server instance router uses the router as gin.Default()
	server.router = router
}

// Starts the http server (define start as a struct method)
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// function to return error as gin object (in case of error required often)
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
