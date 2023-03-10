package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Middleware is a function that wraps an http.Handler to add functionality.
// It takes an http.Handler and returns a new http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain applies a list of middlewares to a http.Handler.
// The middlewares are applied in the order they are provided, meaning the first middleware
// in the list will be the outermost, and the last will be the innermost (closest to the handler).
//