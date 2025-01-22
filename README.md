# URL Shortener Application

## Project Overview

The **URL Shortener Application** is a full-stack web application that allows users to shorten long URLs into compact, shareable links. With features like persistent storage, click tracking, and an intuitive user interface, this project demonstrates modern web development practices and provides a practical solution for managing URLs effectively.

---

## Features

### Core Features

- **Short URL Generation**: Quickly convert long URLs into short, unique links.
- **Redirection**: Redirect users to the original URL when they visit the short link.
- **Persistent Storage**: Store URLs and metadata securely in a PostgreSQL database.

### Advanced Features

- **Analytics & Statistics**: Track link performance with data such as click counts and creation dates.
- **Custom Short Codes**: (Planned) Allow users to create personalized short links.
- **Expiration Dates**: (Planned) Enable links to expire after a predefined period.
- **Authentication**: (Planned) Secure user authentication for managing personalized dashboards.

---

## Technology Stack

### Frontend

- **Framework**: React.js + Typescript
- **Styling**: Tailwind CSS
- **Build Tool**: Vite

### Backend

- **Programming Language**: Go (Golang)
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **ORM**: GORM

### Deployment

- **Frontend**: Deployed on Vercel
- **Backend**: Deployed on Fly.io or Render

---

## Setup Instructions

### Prerequisites

- Node.js and npm
- Go (Golang)
- PostgreSQL
- Docker
- Git

### Backend Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/nickemma/url-shortener.git
   cd url-shortener/server
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Configure the database:

   - Create a PostgreSQL database.
   - Update the `dsn` string in the `initDB` function in `main.go` with your database credentials.

4. Run Docker:

   - Docker-compose up -d
   - The Adminer will start at `http://localhost:8080` for database visualization.

5. Run the backend server:
   ```bash
   go run main.go | make run
   ```
   The server will start at `http://localhost:5000`.

### Frontend Setup

1. Navigate to the frontend directory:

   ```bash
   cd ../Client
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```
   The frontend will run at `http://localhost:5173`.

---

## API Endpoints

### Shorten URL

**POST** `/shorten`

- **Request Body**:
  ```json
  {
    "original_url": "https://example.com"
  }
  ```
- **Response**:
  ```json
  {
    "short_url": "http://localhost:8080/abc123"
  }
  ```

### Redirect to Original URL

**GET** `/:shortCode`

- Redirects to the original URL associated with the `shortCode`.

### Get URL Statistics

**GET** `/stats/:shortCode`

- **Response**:
  ```json
  {
    "original_url": "https://example.com",
    "short_code": "abc123",
    "click_count": 42,
    "created_at": "2025-01-01T12:00:00Z"
  }
  ```

---

## Future Enhancements

- User authentication for managing multiple URLs.
- Analytics dashboard with detailed metrics (e.g., region, browser stats).
- Enhanced security measures to prevent misuse.
- Customizable URL expiration dates.

---

## Contribution Guidelines

We welcome contributions to improve the application. If you have ideas or bug fixes, feel free to:

1. Fork the repository.
2. Create a new branch for your feature/bugfix.
3. Submit a pull request with a detailed description.

---

## License

This project is licensed under the [MIT License](LICENSE).

---

## Acknowledgments

- Inspired by popular URL shortening services.
- Built with love using Go, React.js, Typescript and Tailwind CSS.
