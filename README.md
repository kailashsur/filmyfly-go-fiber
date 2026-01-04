# FilmyFly Go Fiber

A high-performance Go backend for FilmyFly, converted from Node.js/Express to Go/Fiber.

## Features

- **Go Fiber Framework**: Fast HTTP framework built on Fasthttp
- **GORM**: Powerful ORM for PostgreSQL
- **Firebase Authentication**: Admin authentication using Firebase Admin SDK
- **Session Management**: Secure session handling
- **RESTful API**: JSON API for Astro frontend
- **Admin Panel**: EJS-templated admin interface

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Firebase project (for admin authentication)

## Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

3. Install dependencies:

```bash
go mod download
```

4. Run the server:

```bash
go run cmd/server/main.go
```

## Environment Variables

See `.env.example` for all required environment variables.

## API Endpoints

### Public API (for Astro Frontend)

- `GET /api/home` - Homepage data (trending + recent + categories)
- `GET /api/movies` - Paginated movie list
- `GET /api/movies/trending` - Trending movies
- `GET /api/movies/:slug` - Single movie details
- `GET /api/categories` - All categories with counts
- `GET /api/categories/:slug` - Category with movies
- `GET /api/search?q=query` - Search movies
- `GET /api/static-pages/:slug` - Static page content
- `GET /api/astro-settings` - Public settings

### Admin Routes

- `GET /admin/login` - Admin login page
- `POST /admin/login` - Admin login
- `POST /admin/logout` - Admin logout
- `GET /admin` - Admin dashboard (protected)

## Project Structure

```
filmyfly-go-fiber/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and models
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Middleware functions
│   ├── routes/          # Route definitions
│   └── utils/           # Utility functions
├── views/               # EJS templates
├── public/              # Static assets
└── logs/                # Application logs
```

## Development

```bash
# Run with auto-reload (install air first)
air

# Run tests
go test ./...

# Build for production
go build -o filmyfly cmd/server/main.go
```

## License

ISC
