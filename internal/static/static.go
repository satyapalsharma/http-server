package static

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ServeFiles creates an http.Handler that serves static files from the specified root directory.
// It handles URL prefix stripping, provides basic security against directory traversal,
// and adds appropriate cache control headers for production-ready static asset serving.
//
// rootPath: The absolute or relative path to the directory containing static files (e.g., "web/static").
// urlPrefix: The URL path prefix under which these files will be served (e.g., "/static/").
//            It must start with a slash and is recommended to end with one for consistency
//            with http.StripPrefix (e.g., "/static/").
//
// Example usage in your router:
//   router.Handle("/static/", static.ServeFiles("web/static", "/static/"))
func ServeFiles(rootPath, urlPrefix string) http.Handler {
	// --- Configuration Validation and Setup ---

	// 1. Resolve and validate the rootPath.
	// This is a critical startup check; if the static root is invalid, the server should not proceed.
	absRootPath, err := filepath.Abs(rootPath)
	if err != nil {
		log.Fatalf("Static file server: Failed to get absolute path for root '%s': %v", rootPath, err)
	}

	info, err := os.Stat(absRootPath)
	if os.IsNotExist(err) {
		log.Fatalf("Static file server: Root directory '%s' does not exist.", absRootPath)
	}
	if err != nil {
		log.Fatalf("Static file server: Error accessing root directory '%s': %v", absRootPath, err)
	}
	if !info.IsDir() {
		log.Fatalf("Static file server: Root path '%s' is not a directory.", absRootPath)
	}

	// 2. Normalize urlPrefix: ensure it starts and ends with a slash.
	// This makes it consistent for http.StripPrefix and routing.
	if !strings.HasPrefix(urlPrefix, "/") {
		urlPrefix = "/" + urlPrefix
	}
	if !strings.HasSuffix(urlPrefix, "/") {
		urlPrefix = urlPrefix + "/"
	}

	// --- Core File Serving Logic ---

	// 3. Create a file system handler.
	// http.Dir ensures that paths are resolved relative to absRootPath and provides
	// a layer of security against directory traversal by canonicalizing paths.
	fileSystem := http.Dir(absRootPath)
	fileServer := http.FileServer(fileSystem)

	// 4. Wrap the file server with http.StripPrefix.
	// http.StripPrefix removes the urlPrefix from the request URL path
	// before passing it to the fileServer.
	// Example: A request for "/static/css/style.css" with urlPrefix "/static/"
	//          will be stripped to "/css/style.css" for the fileServer to look up.
	strippedFileServer := http.StripPrefix(urlPrefix, fileServer)

	// --- Custom Handler for Enhancements ---

	// 5. Return a custom http.HandlerFunc to