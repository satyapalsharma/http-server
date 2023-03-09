package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-project/http-server/internal/middleware"
	"github.com/your-project/http-server/internal/router"
	"github.com/your-project/http-server/internal/static"
)

// Server represents the HTTP server instance.
// It encapsulates the underlying net/http server, our custom router,
// and provides methods for configuration, starting, and graceful shutdown.
type Server struct {
	httpServer       *http.Server            // The underlying net/http server instance
	router           *router.Router          // Our custom router for handling HTTP requests
	Addr             string                  // The network address the server listens on (e.g., ":8080")
	Logger           *log.Logger             // A logger for server-specific messages and errors
	ShutdownTimeout  time.Duration           // The maximum time allowed for graceful shutdown
	globalMiddleware []middleware.Middleware // A list of middleware functions applied globally to all routes
}

const (
	// defaultShutdownTimeout defines the default duration for graceful server shutdown.
	defaultShutdownTimeout = 15 * time.Second
	// defaultReadTimeout defines the default timeout for reading the entire request, including the body.
	defaultReadTimeout = 5 * time.Second
	// defaultWriteTimeout defines the default timeout for writing the response.
	defaultWriteTimeout = 10 * time.Second
	// defaultIdleTimeout defines the default timeout for keep-alive connections.
	defaultIdleTimeout = 120 * time.Second
)

// New creates and initializes a new Server instance.
// It sets up the server with the given address and logger, and initializes the router.
// If no logger is provided, a default one writing to os.Stdout is used.
func New(addr string, logger *log.Logger) *Server {
	if logger == nil {
		// Fallback to a default logger if none is provided
		logger = log.New(os.Stdout, "[HTTP_SERVER] ", log.LstdFlags|log.Lshortfile)
	}

	r := router.NewRouter()

	s := &Server{
		Addr:             addr,
		Logger:           logger,
		router:           r,
		ShutdownTimeout:  defaultShutdownTimeout,
		globalMiddleware: []middleware.Middleware{}, // Initialize an empty slice for global middleware
	}

	// Initialize the underlying http.Server with sensible defaults.
	// The Handler will be set after all global middleware are applied in the Start method.
	s.httpServer = &http.Server{
		Addr:         s.Addr,
		Handler:      s.router, // Initially set to the raw router, will be wrapped later
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
		ErrorLog:     logger, // Direct http.Server errors to our custom logger
	}

	return s
}

// ApplyGlobalMiddleware adds one or more middleware functions that will be applied
// to the entire router, affecting all registered routes.
// This method should be called before the Server.Start() method.
func (s *Server) ApplyGlobalMiddleware(mw ...middleware.Middleware) {
	s.globalMiddleware = append(s.globalMiddleware, mw...)
}

// RegisterRoutes allows external modules to register their specific routes
// with the server's router. The `registerFunc` typically takes the router
// and adds HTTP handlers for various paths.
func (s *Server) RegisterRoutes(registerFunc func(r *router.Router)) {
	registerFunc(s.router)
}

// ServeStaticFiles configures the server to serve static files from a specified
// file system directory at a given URL prefix.
// For example, ServeStaticFiles("/static/", "web/static") would serve files
// from "web