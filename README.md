# Duolingo Medlab

A Duolingo-style learning platform backend for medical lab science students, built in Go. It serves a JSON API for the learner-facing app and a server-rendered admin panel for managing content (topics, questions, departments).

## Features

- User signup, login, and password reset flows
- Departments and topics for organizing content
- Question bank management via a server-rendered admin panel
- Session-based auth for the admin panel and token-based auth for the API
- Brotli/gzip response compression and request logging middleware

## Tech Stack

- **Language:** Go 1.26
- **Router:** `net/http` (`http.ServeMux`)
- **Database:** PostgreSQL, accessed via [pgx](https://github.com/jackc/pgx)
- **Migrations:** [tern](https://github.com/jackc/tern)
- **Sessions:** `gorilla/securecookie`
- **Validation:** `go-playground/validator`
- **Email:** AWS SES (`aws-sdk-go-v2`)
- **Templates:** Go `html/template` + Tailwind CSS + Alpine.js
- **API docs:** Swagger UI (`swagger-ui/`)

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL (or Docker, to run it via `compose.yml`)
- [tern](https://github.com/jackc/tern) for running migrations
- [air](https://github.com/air-verse/air) (optional, for live reload)

### Setup

1. **Clone the repo and install dependencies**

   ```bash
   go mod download
   ```

2. **Start the database**

   ```bash
   docker compose up -d db
   ```

3. **Configure environment variables**

   Copy `.env` and fill in the values (see [Environment Variables](#environment-variables) below).

4. **Run database migrations**

   ```bash
   tern migrate -m migrations
   ```

5. **Run the server**

   ```bash
   go run main.go
   ```

   Or with live reload via `air`:

   ```bash
   air
   ```

   The server starts on `http://localhost:8080`.

### Environment Variables

| Variable | Description |
| --- | --- |
| `APP_URL` | Base URL of the application |
| `DATABASE_USER` | Postgres username |
| `DATABASE_PASSWORD` | Postgres password |
| `DATABASE_HOST` | Postgres host |
| `DATABASE_PORT` | Postgres port |
| `DATABASE_NAME` | Postgres database name |
| `MAIL_FROM_ADDRESS` | "From" address for outgoing emails (AWS SES) |
| `AWS_REGION` | AWS region for SES |
| `AWS_ACCESS_KEY_ID` | AWS access key |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key |
| `SESSION_HASH_KEY` | Hash key for signing session cookies |
| `SESSION_BLOCK_KEY` | Hex-encoded block key for encrypting session cookies |
| `ADMIN_PASSWORD_HASH` | Bcrypt hash used to seed the admin user (migration-time only) |

## API

The learner-facing API is documented via Swagger UI, served from `/static` assets and the spec at `swagger-ui/api-docs.json`. Key routes:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/api/v1/ping` | Health check |
| `POST` | `/api/v1/auth/signup` | Create a user account |
| `POST` | `/api/v1/auth/login` | Log in and receive a token |
| `POST` | `/api/v1/auth/forgot-password` | Request a password reset |
| `POST` | `/api/v1/auth/reset-password` | Reset password with a token |
| `PATCH` | `/api/v1/users/{id}` | Update a user (auth required) |
| `GET` | `/api/v1/departments` | List departments (auth required) |

## Admin Panel

Server-rendered pages under `/admin`, protected by session auth:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/admin/login` | Admin login page |
| `POST` | `/auth/login` | Admin login submit |
| `POST` | `/admin/logout` | Log out |
| `GET` | `/admin/dashboard` | Dashboard |
| `POST` | `/admin/topics` | Create a topic |
| `GET` | `/admin/questions` | List questions |
| `GET` | `/admin/questions/create` | New question form |
| `POST` | `/admin/questions` | Create a question |
| `GET` | `/admin/questions/{id}` | View a question |
| `POST` | `/admin/questions/{id}/delete` | Delete a question |

## Project Structure

```
handlers/       HTTP handlers (API + admin)
pkg/
  middleware/   Auth, session, compression, content-type middleware
  models/       Database models
  request/      Request DTOs and validation
  response/     Response DTOs
  mailer/       Email sending (AWS SES)
migrations/     SQL migrations (run via tern)
templates/      Admin panel HTML templates
static/         CSS/JS assets served at /static
swagger-ui/     API documentation
```

## Testing

```bash
go test ./...
```

## Deployment

The app is containerized via `Dockerfile` and deployed with `compose.prod.yml`, which runs behind [Traefik](https://traefik.io/) with automatic TLS.

```bash
docker compose -f compose.prod.yml up -d --build
```
