# GoLink.smol

A modern, high-performance URL shortener built with Go, HTMX, and Redis. GoLink.smol uses SHA-256 hashing with collision handling to generate unique, shortened URLs while ensuring fast lookups and redirects.

## Features

- **Fast URL Shortening**: Generates short, unique URLs using SHA-256 hashing
- **Collision Detection**: Automatically handles hash collisions by regenerating keys
- **Modern UI**: Built with HTMX for seamless, JavaScript-free interactions
- **Redis Backend**: High-performance Redis storage with quick URL lookups
- **Link Validation**: Verifies destination URLs before shortening
- **Graceful Shutdown**: Handles server termination gracefully
- **15-day Expiration**: URLs automatically expire after 15 days

## Tech Stack

- **Backend**: Go
- **Frontend**: HTMX
- **Database**: Redis
- **Template Engine**: Go's built-in HTML templating

## Prerequisites

- Go 1.22 or higher
- Docker (for Redis)

## Installation and Setup

1. Clone the repository:

```bash
git clone https://github.com/yourusername/golink.smol.git
cd golink.smol
```

2. Create a `.env` file in the project root:

```env
# port for http server
PORT=8080
# conn string for redis
CONN_STR=redis://:log123@localhost:6379
```

3. Start Redis with Docker:

```bash
docker run --name redis-server -p 6379:6379 -d redis
```

4. Install dependencies:

```bash
go mod vendor
go mod tidy
```

5. Build and run binary:

```bash
go build -o bin/smol.exe
./bin/smol.exe
```

The server will be available at `http://localhost:8080`

## Project Structure

```
.
├── main.go             # Entry point, server initialization
├── storage.go          # Redis store implementation
├── handlers.go         # HTTP handlers
├── api.go              # Template and route setup
├── utils.go            # Utility functions
├── .env                # Environment configuration
├── templates/          # HTML templates
└── static/             # Static assets
```

## Architecture

### URL Shortening Process

1. **URL Submission**: User submits a URL through the HTMX-powered form
2. **Validation**: Server validates the URL by attempting to connect to it
3. **Hash Generation**:
   - Creates a SHA-256 hash of the URL with a random byte for uniqueness
   - Uses first 6 bytes of hash for the short URL
   - Checks for collisions and regenerates if necessary
4. **Storage**: Stores the URL mapping in Redis with a 15-day expiration

### Key Components

- **RedisStore**: Handles all Redis operations and URL encoding
- **GoLinkServer**: Main server implementation with routing and handlers
- **Templates**: Server-side rendered HTML templates
- **Static Files**: Serves static assets from the `/static` directory

## API Endpoints

- `POST /make-it-smol`: Shortens a URL
- `GET /{id}`: Redirects to the original URL
- `GET /`: Home page
- `GET /not-found`: 404 error page

## Error Handling

- Invalid URLs return a 301 redirect to `/not-found`
- Redis connection errors trigger server shutdown
- Graceful shutdown with 30-second timeout
- Link validation with 7-second timeout

## Security Features

- Random byte addition to prevent hash collisions
- URL validation before shortening
- Limited URL key length (6 bytes)
- Automatic URL expiration

## Development

To contribute to GoLink.smol:

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Submit a pull request

## License

GPL-3.0 License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
