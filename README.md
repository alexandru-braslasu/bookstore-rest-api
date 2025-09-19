# Bookstore API

A RESTful API service built with Go and Gin framework for bookstore management, featuring user authentication, inventory control, and personalized book collections.

## Overview

This application provides a comprehensive backend solution for bookstore operations, including book catalog management, user authentication with multiple providers, and individual user collections. The system supports both traditional email/password authentication and Google OAuth integration.

## Features

- **Book Management**: Complete CRUD operations for book inventory
- **User Authentication**: JWT-based authentication with Google OAuth support  
- **Role-Based Access Control**: Administrative and standard user permissions
- **Personal Collections**: Individual user book libraries with quantity tracking
- **Search Functionality**: Book discovery by title and author

## Technology Stack

- **Runtime**: Go 1.19+
- **Web Framework**: Gin HTTP framework
- **Database**: SQLite3 (development), production-ready for PostgreSQL/MySQL
- **Authentication**: JWT tokens, Google OAuth 2.0
- **Security**: bcrypt password hashing

### Dependencies

```
github.com/gin-gonic/gin          # HTTP web framework
github.com/golang-jwt/jwt/v5      # JWT implementation
github.com/markbates/goth         # OAuth provider integration
golang.org/x/crypto/bcrypt        # Password hashing
github.com/mattn/go-sqlite3       # SQLite database driver
github.com/joho/godotenv          # Environment variable management
```

### Setup Instructions

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd bookstore-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Configure environment variables**
   
   Create a `.env` file in the project root:
   ```env
   SECRET_KEY=your-jwt-secret-key
   CLIENT_ID=your-google-oauth-client-id
   CLIENT_SECRET=your-google-oauth-client-secret
   ```

4. **Google OAuth Setup**
   - Navigate to [Google Cloud Console](https://console.cloud.google.com/)
   - Create or select a project
   - Enable the Google+ API
   - Create OAuth 2.0 credentials
   - Configure authorized redirect URI: `http://localhost:8080/bookstore/users/auth`

5. **Start the application**
   ```bash
   go run main.go
   ```

The server will be available at `http://localhost:8080`

## Database Schema

The application automatically initializes the following SQLite tables:

### Books
```sql
CREATE TABLE books (
    title TEXT PRIMARY KEY,
    author TEXT NOT NULL,
    description TEXT,
    nrSamples INTEGER NOT NULL
);
```

### Users
```sql
CREATE TABLE users (
    name TEXT NOT NULL,
    email TEXT PRIMARY KEY,
    password TEXT NOT NULL,
    isAdmin BOOLEAN
);
```

### Saved Books
```sql
CREATE TABLE saved_books (
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    description TEXT NOT NULL,
    email TEXT NOT NULL,
    nrSamples INTEGER NOT NULL,
    PRIMARY KEY (title, author, email),
    FOREIGN KEY (email) REFERENCES users(email),
    FOREIGN KEY (title, author) REFERENCES books(title, author)
);
```

## API Documentation

### Public Endpoints

#### Book Catalog
- `GET /bookstore/books` - Retrieve all available books
- `GET /bookstore/books/findByTitle/:title` - Find book by exact title match
- `GET /bookstore/books/findByAuthor/:author` - Find books by author (**obs!: supports partial matching**)

#### Authentication
- `POST /bookstore/users/signup` - User registration
- `POST /bookstore/users/login` - Email/password authentication
- `GET /bookstore/users/loginWithGoogle` - Initialize Google OAuth flow
- `GET /bookstore/users/auth` - Google OAuth callback handler

### Protected Endpoints

All protected endpoints require authentication via the `Authorization` header containing a valid JWT token.

#### Administrative Operations
- `POST /bookstore/books` - Create new book entry (Admin only)

#### User Operations
- `PUT /bookstore/users/saveBook` - Add book to personal collection or add additional samples to one already saved
- `GET /bookstore/users/savedBooks` - Retrieve user's saved books

## Request/Response Examples

### User Registration
```bash
curl -X POST http://localhost:8080/bookstore/users/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Legend",
    "email": "john.legend@example.com",
    "password": "givemehope"
  }'
```

**Response:**
```json
{
  "message": "User added successfully!"
}
```

### Authentication
```bash
curl -X POST http://localhost:8080/bookstore/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.legend@example.com",
    "password": "givemehope"
  }'
```

**Response:**
```json
{
  "message": "Login successful!",
  "token": "generated_token"
}
```

### Book Creation (Admin)
```bash
curl -X POST http://localhost:8080/bookstore/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{
    "title": "Klara and the Sun",
    "author": "Kazuo Ishiguro",
    "description": "From her place in the store, Klara, an Artificial Friend ...",
    "nrSamples": 25
  }'
```

### Save Book to Collection
```bash
curl -X PUT http://localhost:8080/bookstore/users/saveBook \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{
    "title": "Klara and the Sun",
    "author": "Kazuo Ishiguro",
    "nrSamples": 1
  }'
```

## Authentication Flow

#### Standard Authentication
1. User registration creates account with hashed password
2. Login endpoint validates credentials and returns JWT token
3. Protected routes require valid JWT token in Authorization header
4. Tokens expire after 2 hours and must be renewed

#### Google OAuth Flow
1. Client redirects to `/bookstore/users/loginWithGoogle`
2. User completes Google authentication
3. System automatically add new special user to database if it does not exist yet
4. Authentication callback handled at `/bookstore/users/auth`

## Security Implementation

### Password Security
- All passwords are hashed using bcrypt with cost factor 14
- Plain-text passwords are never stored in the database
- Password validation occurs during login through hash comparison

### JWT Token Management
- Tokens are signed using HMAC SHA-256 algorithm
- Each token contains user email, admin status, and expiration time
- Token validation includes signature verification and expiration checking

### Authorization Middleware
- All protected routes pass through authentication middleware
- Supports both JWT and Google OAuth token validation
- User context (email, admin status) is injected into request context