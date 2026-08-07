# Klinme

> AI-powered data cleaning SaaS for agencies and businesses.

Klinme helps agencies and companies clean messy CSV and Excel files using a combination of rule-based logic and AI — fast, simple, and affordable.

---

## Table of Contents

- [Overview](#overview)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
- [Database](#database)
- [Authentication](#authentication)
- [File Upload](#file-upload)
- [Roadmap](#roadmap)

---

## Overview

Klinme is a SaaS platform that allows agencies and companies to upload CSV/Excel files and get them cleaned automatically. The cleaning pipeline combines rule-based logic (trim whitespace, fix date formats, remove duplicates) with AI to handle ambiguous cases.

**Key Features:**

- Upload CSV and Excel files (up to 10MB)
- Automatic rule-based + AI-powered data cleaning
- User authentication and session management via Clerk
- File storage via Azure Blob Storage
- Usage-based pricing with subscription plans
- Per-user cleaning history and job tracking

---

## Tech Stack

| Layer          | Technology                |
| -------------- | ------------------------- |
| Backend        | Go (Golang) + Gin         |
| Database       | PostgreSQL (Neon)         |
| Authentication | Clerk                     |
| File Storage   | Azure Blob Storage        |
| Payments       | Paystack (coming soon)      |
| Deployment     | Azure (coming soon)       |

---

## Project Structure

```
klinme-api/
├── main.go                  # Entry point
├── .env                     # Environment variables (never commit)
├── .gitignore
├── go.mod
├── go.sum
├── db/
│   └── db.go                # Database connection
├── handlers/
│   ├── ping.go              # Health check
│   ├── user.go              # User handlers
│   ├── file.go              # File upload handler
│   └── webhook.go           # Clerk webhook handler
├── middleware/
│   └── auth.go              # Clerk JWT auth middleware
├── models/
│   ├── user.go              # User model
│   ├── file.go              # File model
│   └── cleanjob.go          # CleanJob model
├── routes/
│   └── routes.go            # All routes
├── storage/
│   └── azure.go             # Azure Blob Storage client
└── migrations/
    └── 001_init.sql         # Database schema
```

---

## Getting Started

### Prerequisites

- Go 1.22+
- PostgreSQL 16+
- A [Neon](https://neon.tech) account
- A [Clerk](https://clerk.com) account
- An [Azure](https://portal.azure.com) account

### Installation

```bash
# Clone the repository
git clone https://github.com/Samueljr-web/klinme-api.git
cd klinme-api

# Install dependencies
go mod tidy

# Copy environment variables
cp .env.example .env
# Fill in your values in .env

# Run database migrations
psql "YOUR_DATABASE_URL" -f migrations/001_init.sql

# Start development server
air
```

---

## Environment Variables

Create a `.env` file in the root of the project:

```env
# Database
DATABASE_URL=postgresql://user:password@host/dbname?sslmode=require

# Clerk Authentication
CLERK_SECRET_KEY=sk_test_...
CLERK_PUBLISHABLE_KEY=pk_test_...

# Azure Blob Storage
AZURE_STORAGE_ACCOUNT=your_storage_account
AZURE_STORAGE_KEY=your_storage_key
AZURE_STORAGE_CONNECTION_STRING=your_connection_string
AZURE_RAW_CONTAINER=raw
AZURE_CLEANED_CONTAINER=cleaned
```

> Never commit your `.env` file. It is listed in `.gitignore`.

---

## API Endpoints

### Public Routes

| Method | Endpoint              | Description                     |
| ------ | --------------------- | ------------------------------- |
| GET    | `/api/ping`           | Health check                    |
| POST   | `/api/webhooks/clerk` | Clerk webhook for user creation |

### Protected Routes (require Bearer token)

| Method | Endpoint            | Description                |
| ------ | ------------------- | -------------------------- |
| GET    | `/api/users/:id`    | Get user by ID             |
| POST   | `/api/files/upload` | Upload a CSV or Excel file |

### Authentication

All protected routes require an `Authorization` header:

```
Authorization: Bearer YOUR_CLERK_SESSION_TOKEN
```

---

## Database

### Schema

**users**
| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| clerk_user_id | TEXT | Clerk user ID |
| email | TEXT | User email |
| plan | TEXT | free, starter, growth |
| cleans_used | INT | Number of cleans used |
| created_at | TIMESTAMP | Account creation date |

**files**
| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| user_id | UUID | References users.id |
| file_name | TEXT | Original file name |
| file_size | BIGINT | File size in bytes |
| file_type | TEXT | .csv, .xlsx, .xls |
| status | ENUM | pending, processing, done, failed |
| created_at | TIMESTAMP | Upload date |

**clean_jobs**
| Column | Type | Description |
|---|---|---|
| id | UUID | Primary key |
| user_id | UUID | References users.id |
| file_id | UUID | References files.id |
| status | ENUM | pending, processing, done, failed |
| rows_processed | INT | Total rows processed |
| rows_cleaned | INT | Rows that were cleaned |
| error_message | TEXT | Error details if failed |
| created_at | TIMESTAMP | Job creation date |
| completed_at | TIMESTAMP | Job completion date |

---

## Authentication

Klinme uses [Clerk](https://clerk.com) for authentication.

- User signup and login is handled entirely by Clerk on the frontend
- On successful signup, Clerk fires a `user.created` webhook to `/api/webhooks/clerk`
- Your Go backend receives the webhook and creates the user in PostgreSQL
- All subsequent API requests are authenticated via Clerk JWT tokens verified in the `AuthMiddleware`

---

## File Upload

- Supported formats: `.csv`, `.xlsx`, `.xls`
- Maximum file size: **10MB**
- Raw files are stored in Azure Blob Storage under the `raw/` container
- Cleaned files will be stored under the `cleaned/` container
- Every upload creates a `files` record and a `clean_jobs` record in PostgreSQL

---

## Roadmap

- [x] Project setup (Go + Gin)
- [x] PostgreSQL database + migrations
- [x] Azure Blob Storage integration
- [x] Clerk authentication + webhook
- [x] File upload endpoint
- [x] Rule-based data cleaning logic
- [x] AI-powered cleaning integration
- [ ] Paystack subscription + usage limits
- [ ] Production deployment
