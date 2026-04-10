# 🧠 TaskFlow Backend — Production-Ready Go API

A clean, production-quality backend built in Go for managing projects and tasks with authentication, authorization, and structured architecture.

This project is designed to demonstrate **real backend engineering practices** — not just CRUD APIs.

---

# 🚀 Overview

TaskFlow is a backend system that allows users to:

* Register and authenticate using JWT
* Create and manage projects
* Create, update, and assign tasks within projects
* Filter and query tasks efficiently
* Enforce strict authorization rules

The system is built with a strong focus on:

* Clean architecture
* Separation of concerns
* Production-ready practices
* Observability and maintainability

---

# 🏗️ Architecture

The project follows a layered architecture:

```text
Handler → Service → Repository → Database
```

### 🔹 Handler Layer

* Handles HTTP requests/responses
* Input validation
* Status code management

### 🔹 Service Layer

* Business logic
* Authorization rules
* Data validation

### 🔹 Repository Layer

* Database interaction
* Raw SQL queries
* No ORM (as required)

---

# 🔐 Authentication

* JWT-based authentication
* Token contains:

  * `user_id`
  * `email`
* Token expiry: **24 hours**
* Password hashing using **bcrypt (cost 12)**
* All protected routes require:

```http
Authorization: Bearer <token>
```

---

# 🛡️ Authorization Rules

### Projects

* Only **owner** can update/delete a project

### Tasks

* Only **project owner** can update/delete tasks
  *(simplified for this implementation)*

> Note: In a production system, task-level ownership (creator) can be added for finer control.

---

# 📦 Tech Stack

* **Language:** Go
* **Framework:** Gin
* **Database:** PostgreSQL
* **Migrations:** golang-migrate
* **Auth:** JWT
* **Logging:** slog (structured logging)
* **Containerization:** Docker + Docker Compose

---

# ⚙️ Getting Started

## 1. Clone repository

```bash
git clone <your-repo-url>
cd taskflow-backend
```

---

## 2. Setup environment

```bash
cp .env.example .env
```

---

## 3. Run the application

```bash
docker compose up --build
```

---

## 4. Server will be available at

```text
http://localhost:<API_PORT>
```

---

# 🔑 Test Credentials

```text
email: test@example.com
password: password123
```

---

# 📚 API Reference

---

## 🔐 Auth

### Register

```http
POST /auth/register
```

### Login

```http
POST /auth/login
```

---

## 📁 Projects

### Get all projects

```http
GET /projects
```

### Create project

```http
POST /projects
```

### Update project

```http
PATCH /projects/:id
```

### Delete project

```http
DELETE /projects/:id
```

---

## 🧩 Tasks

### Get tasks (with filters)

```http
GET /projects/:id/tasks?status=&assignee=
```

### Create task

```http
POST /projects/:id/tasks
```

### Update task

```http
PATCH /tasks/:id
```

### Delete task

```http
DELETE /tasks/:id
```

---

## 📊 Project Stats (Bonus)

### Get aggregated stats

```http
GET /projects/:id/stats
```

Returns:

* Task count by status
* Task count by assignee

---

# 🧠 Key Design Decisions

---

## 1. Clean Architecture

Separated layers ensure:

* Maintainability
* Testability
* Scalability

---

## 2. No ORM

Used raw SQL because:

* Better control over queries
* Matches assignment requirement
* Improves performance awareness

---

## 3. Simplified Task Authorization

* Only project owners can modify tasks

**Tradeoff:**

* Simpler logic
* Easier to reason about

**Future improvement:**

* Add `creator_id` for finer access control

---

## 4. Structured Logging (slog)

* Logs include method, path, status, latency
* Designed for production observability

---

## 5. UUID-based IDs

* Ensures global uniqueness
* Avoids predictable IDs

---

# 🧪 Testing

A Postman/Bruno collection is included covering:

* Auth flows
* Project operations
* Task operations
* Authorization edge cases

---

# 🐳 Docker Setup

* PostgreSQL container with healthcheck
* API container depends on DB readiness
* Migrations run automatically on startup

---

# 📈 What I Would Improve With More Time

* Add pagination for large datasets
* Introduce role-based access control (RBAC)
* Add caching layer (Redis)
* Add rate limiting and request throttling
* Write full integration tests
* Add CI/CD pipeline
* Add search functionality

---

# 💡 Final Thoughts

This project reflects how I approach backend systems:

* Start simple, but design for scale
* Prioritize correctness and clarity
* Enforce rules at the right layer
* Build with production in mind

---

> “A good backend is not just about endpoints — it's about enforcing rules, maintaining consistency, and being predictable under all conditions.”
