```go
package router

import (
	"log"
	"net/http"
	"strings"
	"sync"
)

// Middleware is a function that takes an http.Handler and returns an http.Handler.
// It's used to wrap handlers with additional logic (e.g., logging, authentication).
type Middleware func(http.Handler) http.Handler

// Router is a custom HTTP router that supports middleware and method-specific routing.
// It implements the http.Handler interface, allowing it to be used directly with http.Server.
type Router struct {
	// routes stores handlers for specific HTTP methods and paths.
	// The key is a string combining the HTTP method and the path (e.g., "GET /