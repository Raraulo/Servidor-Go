<p align="center">
  <img src="https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_Blue.png" width="220" alt="Go Logo">
</p>

<h1 align="center">Servidor Go</h1>

<p align="center">
  RESTful API developed in Go for an e-commerce platform, providing authentication, catalog management, inventory control, sales, reporting and seamless integration with an Angular frontend.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21-00ADD8?logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/Gorilla_Mux-Router-000000" alt="Gorilla Mux">
  <img src="https://img.shields.io/badge/GORM-ORM-00A98F" alt="GORM">
  <img src="https://img.shields.io/badge/PostgreSQL-Database-4169E1?logo=postgresql&logoColor=white" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/JWT-Authentication-000000?logo=jsonwebtokens" alt="JWT">
  <img src="https://img.shields.io/badge/CORS-Enabled-28A745" alt="CORS">
  <img src="https://img.shields.io/badge/Render-Deployment-46E3B7?logo=render&logoColor=white" alt="Render">
</p>

---

## Table of Contents

- Overview
- Features
- Technology Stack
- Architecture
- API Modules
- Project Structure
- Installation
- Environment Variables
- Development
- Deployment
- Security
- Available Commands
- Project Status
- Author

---

# Overview

Servidor Go is a RESTful backend developed with **Go** to support an e-commerce platform for pet products.

The API provides secure authentication, product management, inventory control, purchase processing, invoice generation and administrative reporting. It serves as the primary backend for an **Angular frontend**, exposing REST endpoints protected with JWT authentication.

The application is designed for local development as well as cloud deployment using **PostgreSQL** and **Render**.

---

# Features

## Authentication

- JWT authentication
- User registration
- Secure login
- Role-based authorization

## Catalog Management

- Product management
- Category management
- Pet type management
- Inventory control

## Sales

- Purchase processing
- Invoice generation
- Purchase history
- Stock validation

## Administrative Reports

- Sales reports
- Product statistics
- Inventory reports
- Customer reports

## API

- RESTful architecture
- JSON responses
- CORS enabled
- Angular frontend integration

---

# Technology Stack

| Category | Technology |
|-----------|------------|
| Language | Go 1.21 |
| Router | Gorilla Mux |
| ORM | GORM |
| Database | PostgreSQL |
| Authentication | JWT |
| Security | CORS |
| Deployment | Render |

---

# Architecture

The backend is organized into independent modules responsible for different areas of the application.

- Authentication
- Users
- Products
- Categories
- Pet Types
- Sales
- Invoices
- Reports
- Database

---

# API Modules

| Module | Description |
|----------|-------------|
| Authentication | JWT authentication and authorization |
| Users | User management |
| Products | Product catalog |
| Categories | Product categories |
| Animals | Pet type management |
| Sales | Purchase processing |
| Invoices | Invoice generation |
| Reports | Administrative statistics |

---

# Project Structure

```text
Servidor-Go/
│
├── db/
│   └── baseDb.go
│
├── dominios/
│   ├── animal.go
│   ├── categoria.go
│   ├── factura.go
│   ├── persona.go
│   ├── producto.go
│   └── venta.go
│
├── handlers/
│   ├── seguridad/
│   └── utiles/
│
├── rutas/
│   ├── rt_animal.go
│   ├── rt_fact.go
│   ├── rt_prdcto.go
│   ├── rt_reportes.go
│   └── rt_usro.go
│
├── token/
│   ├── auth.go
│   ├── config.go
│   ├── jwt.go
│   └── middleware.go
│
├── main.go
├── go.mod
├── go.sum
└── gdo.sql
```

---

# Installation

Clone the repository.

```bash
git clone <repository-url>
cd Servidor-Go
```

Install dependencies.

```bash
go mod tidy
```

Run the server.

```bash
go run main.go
```

The application runs locally at:

```text
http://localhost:4000
```

---

# Environment Variables

Create the required environment variables before running the application in production.

```env
DATABASE_URL=

PORT=

JWT_SECRET=
```

Sensitive credentials should never be committed to the repository.

---

# Development

Useful development commands.

Run the server.

```bash
go run main.go
```

Download project dependencies.

```bash
go mod tidy
```

Build the application.

```bash
go build
```

Generate a custom executable.

```bash
go build -o servidor
```

Run the compiled executable.

Linux/macOS

```bash
./servidor
```

Windows

```bash
servidor.exe
```

---

# Deployment

The application is optimized for deployment on **Render**.

Recommended configuration:

| Setting | Value |
|----------|-------|
| Runtime | Go |
| Build Command | `go build -o servidor` |
| Start Command | `./servidor` |

Required environment variables:

```env
DATABASE_URL=
PORT=
JWT_SECRET=
```

---

# Security

The API uses JWT authentication to secure protected endpoints.

Access levels include:

- Public endpoints
- Authenticated user endpoints
- Administrative endpoints

Security recommendations:

- Never commit environment files.
- Never expose database credentials.
- Keep JWT secrets private.
- Restrict administrative endpoints.
- Enable HTTPS in production.

---

# Available Commands

| Command | Description |
|----------|-------------|
| go mod tidy | Install project dependencies |
| go run main.go | Start development server |
| go build | Compile the application |
| go build -o servidor | Generate executable |
| ./servidor | Run compiled application (Linux/macOS) |
| servidor.exe | Run compiled application (Windows) |

---

# Project Status

The application is production-ready and includes:

- RESTful API
- PostgreSQL integration
- Automatic database migrations
- JWT authentication
- Inventory management
- Purchase processing
- Administrative reporting
- Role-based authorization
- Angular frontend integration
- Cloud deployment support

---

# Author

RESTful backend developed with **Go**, **Gorilla Mux**, **GORM** and **PostgreSQL**, designed to power an Angular-based e-commerce platform through secure and scalable API services.
