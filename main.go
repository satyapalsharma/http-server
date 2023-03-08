package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"myproject/internal/api"
	"myproject/internal/middleware"
	"myproject/internal/router"
	"myproject/internal/server"
	"myproject/internal/static"
)

const (
	// defaultPort is the fallback port if not specified in environment variables.
	defaultPort = "8080"
	// shutdownTimeout is the maximum time allowed for graceful server shutdown.
	