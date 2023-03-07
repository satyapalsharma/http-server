# Lightweight Go HTTP Server

A lightweight and modular HTTP server built with Go's `net/http` package, featuring robust routing, extensible middleware, static file serving, and JSON API capabilities. Designed for simplicity, performance, and ease of extension.

## Features

*   **Custom Router**: Define routes for different HTTP methods (`GET`, `POST`, `PUT`, `DELETE`).
*   **Middleware Support**: Apply reusable logic (e.g., logging, authentication, CORS) to requests.
*   **Static File Serving**: Efficiently serve HTML, CSS, JavaScript, images, and other static assets.
*   **JSON API Endpoints**: Easily build RESTful APIs with structured JSON request and response handling.
*   **Modular Design**: Clear separation of concerns for server setup, routing, middleware, and handlers, promoting maintainability.
*   **Environment Configuration**: Utilizes `.env` files for flexible configuration.

## Tech Stack

*   **Go**: The primary language for the server logic.
*   **`net/http`**: Go's standard library for HTTP server functionality.
*   **`github.com/joho/godotenv`**: For loading environment variables from `.env` files.

## Project Structure

The project is organized to promote modularity and maintainability:

```
.
├── .env.example                 # Example environment variables
├── .gitignore                   # Git ignore rules
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── main.go                      # Main application entry point
├── README.md                    # This file
├── internal/                    # Internal packages for server logic
│   ├── api/                     # Handlers for JSON API endpoints
│   │   └── handlers.go
│   ├── middleware/              # Reusable HTTP middleware functions
│   │   └── middleware.go
│   ├── router/                  # Custom router implementation
│   │   └── router.go
│   ├── server/                  # Core server setup and configuration
│   │   └── server.go
│   └── static/                  # Handler for serving static files
│       └── static.go
└── web/                         # Web assets
    └── static/                  # Static files to be served
        ├── css/
        │   └── style.css
        └── index.html
```

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

*   Go (version 1.18 or higher)
    *   [Download and Install Go](https://golang.org/doc/install)

### Installation and Running

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/go-http-server.git
    cd go-http-server
    ```

2.  **Set up environment variables:**
    Copy the `.env.example` file to `.env` and modify it as needed.

    ```bash
    cp .env.example .env
    ```

    The `.env` file should look something like this:

    ```env
    PORT=8080
    STATIC_DIR=./web/static
    ```

3.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

4.  **Run the server:**

    ```bash
    go run main.go
    ```

    The server will start on the port specified in your `.env` file (default: `8080`). You should see output similar to:

    ```
    Loading environment variables from .env
    Server starting on :8080
    ```

## Usage

Once the server is running, you can access it through your web browser or API client.

### Static Files

*   Open your browser to `http://localhost:8080/`
    *   This will serve the `index.html` file from `web/static`.
*   Access static CSS: `http://localhost:8080/css/style.css`

### API Endpoints

The server provides a few example API endpoints:

*   **`GET /api/hello`**
    *   Returns a simple JSON greeting.
    *   Example response: `{"message": "Hello from Go API!"}`

*   **`POST /api/data`**
    *   Accepts a JSON payload and returns a processed response.
    *   **Request Example (using `curl`):**
        ```bash
        curl -X POST -H "Content-Type: application/json" \
             -d '{"name": "Alice", "age": 30}' \
             http://localhost:8080/api/data
        ```
    *   **Response Example:**
        ```json
        {
            "received_name": "Alice",
            "received_age": 30,
            "status": "Data processed successfully"
        }
        ```

*   **`GET /api/protected`**
    *   This endpoint demonstrates middleware. It's protected by a simple `AuthMiddleware` that checks for an `X-API-KEY` header.
    *   **Without `X-API-KEY` header:**
        ```bash
        curl http://localhost:8080/api/protected
        ```
        *   Response: `{"error": "Unauthorized"}` (Status: 401 Unauthorized)
    *   **With `X-API-KEY` header:**
        ```bash
        curl -H "X-API-KEY: mysecretkey" http://localhost:8080/api/protected
        ```
        *   Response: `{"message": "Welcome, authorized user!"}`

### Middleware in Action

Observe your server console when making requests. The `LoggingMiddleware` will print details for each incoming request, demonstrating how middleware intercepts and processes requests before they reach the final handler.

```
[INFO] 2023-10-27 10:00:00 GET /api/hello
[INFO] 2023-10-27 10:00:01 POST /api/data
[INFO] 2023-10-27 10:00:02 GET /api/protected
```

## Configuration

The server uses environment variables for configuration, loaded from a `.env` file.

*   **`PORT`**: The port on which the server will listen (e.g., `8080`).
*   **`STATIC_DIR`**: The path to the directory containing static files to be served (e.g., `./web/static`).

Ensure your `.env` file is correctly configured before starting the server.

## Contributing

Contributions are welcome! If you have suggestions for improvements, bug fixes, or new features, please feel free to:

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature/your-feature-name`).
3.  Make your changes and commit them (`git commit -m 'Add new feature'`).
4.  Push to the branch (`git push origin feature/your-feature-name`).
5.  Open a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.