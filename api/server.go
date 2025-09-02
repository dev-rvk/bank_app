package api

import (
	db "github.com/devrvk/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// server struct contains a server router and store
type Server struct {
	store db.Store
	router *gin.Engine
}

// constructor for the server struct with routes initialized
func NewServer (store db.Store) *Server {

	// custom validator for currency
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok{
		v.RegisterValidation("currency", validCurrency)
	}

	// takes store from arguments and uses gin Default router
	server := &Server{store: store}
	router := gin.Default()

	// define route
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)

	// server instance router uses the router as gin.Default()
	server.router = router

	// returns server object of type Server contains store and router
	return server
}

// Starts the http server (define start as a struct method)
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// function to return error as gin object (in case of error required often)
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}