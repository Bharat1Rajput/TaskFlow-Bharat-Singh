# TaskFlow Backend

A small but production-minded backend built in Go to manage projects and tasks with authentication and proper access control.

I focused on writing clean, readable code and enforcing rules at the right layer rather than just exposing endpoints.

---

## Overview

This service allows users to:

* Register and login using JWT
* Create and manage projects
* Create and manage tasks inside projects
* Filter tasks by status and assignee
* Enforce ownership rules on updates and deletes

**Tech stack:**

* Go (Gin)
* PostgreSQL
* JWT for auth
* golang-migrate for migrations
* Docker + docker compose

---

## Architecture Decisions

I followed a simple layered structure:

```text
handler → service → repository → database
```

* **Handlers** deal with HTTP (validation, response codes)
* **Services** contain business logic and authorization rules
* **Repositories** handle raw SQL queries

### Why this structure?

It keeps things predictable. Each layer has one job, which makes the code easier to reason about and debug.

### Tradeoffs I made

* I **restricted task updates/deletes to project owners only**
  → I intentionally skipped `creator_id` to keep authorization simple within the assignment scope

* I used **raw SQL instead of an ORM**
  → more control and aligns with the assignment, but slightly more verbose

* I didn’t over-engineer (no caching, no RBAC, no microservices)
  → focused on correctness and clarity first

---

## Running Locally

Assuming you only have Docker installed:

```bash
git clone https://github.com/Bharat1Rajput/taskflow-Bharat-Singh
cd taskflow-Bharat-Singh

cp .env.example .env

docker compose up --build
```

Once everything is up, the API will be available at:

```text
http://localhost:<API_PORT>
```

---

## Running Migrations

Migrations run automatically when the application starts using `golang-migrate`.

No manual commands are needed.

---

## Test Credentials

You can log in immediately using the seeded user:

```text
Email:    test@example.com
Password: Bharat123
```

---

## API Reference

### Auth

* `POST /auth/register`
* `POST /auth/login`

---

### Projects

* `GET /projects`
* `POST /projects`
* `PATCH /projects/:id`
* `DELETE /projects/:id`

---

### Tasks

* `GET /projects/:id/tasks?status=&assignee=`
* `POST /projects/:id/tasks`
* `PATCH /tasks/:id`
* `DELETE /tasks/:id`

---

### Example Flow

1. Register or login → get JWT
2. Create a project
3. Add tasks to the project
4. Use filters (`status`, `assignee`) to query tasks

All non-auth routes require:

```http
Authorization: Bearer <token>
```

---

## What I’d Do With More Time

* Add pagination to project and task listing
* Introduce finer-grained permissions (task creator vs owner)
* Add integration tests
* Add rate limiting and request validation improvements
* Improve error handling consistency across all endpoints

---

## Final Note

This project is intentionally kept simple in terms of scope, but I tried to keep the internals clean and realistic.

The goal was not just to “make it work”, but to make it understandable and maintainable.
