# GDU-Go-Init

A small Go backend project demonstrating CRUD operations. Built with Gin, GORM, and MariaDB.

## Stack

- **Go** with [Gin](https://github.com/gin-gonic/gin) — HTTP router
- **GORM** — ORM for database access
- **MariaDB** — relational database
- **JWT** — authentication tokens

---

## Run with Docker Compose

Start the database:

```bash
docker-compose up -d
```

Then run the Go server:

```bash
go run main.go
```

The API will be available at `http://localhost:8080`.

> MariaDB runs on port `3306`. Data is persisted in `./data/mysql/`.

---

## API Endpoints

### Users

| Method | Endpoint           | Description              |
|--------|--------------------|--------------------------|
| POST   | `/users/register`  | Create a new account     |
| POST   | `/users/login`     | Login and get a JWT token |

**Register**
```json
POST /users/register
{
  "email": "you@example.com",
  "password": "secret123"
}
```

**Login**
```json
POST /users/login
{
  "email": "you@example.com",
  "password": "secret123"
}
```
Returns a JWT token valid for 2 hours.

---

### Projects

| Method | Endpoint          | Description            |
|--------|-------------------|------------------------|
| GET    | `/projects/`      | List all projects      |
| GET    | `/projects/:id`   | Get one project        |
| POST   | `/projects/`      | Create a project       |
| PUT    | `/projects/:id`   | Update a project       |
| DELETE | `/projects/:id`   | Delete a project       |

**Create a project**
```json
POST /projects/
{
  "name": "My Project",
  "description": "A cool project",
  "image": "https://example.com/image.png",
  "skills": ["Go", "Docker", "SQL"]
}
```

**Update a project** (partial update — only send the fields you want to change)
```json
PUT /projects/1
{
  "name": "Updated name"
}
```

---

## Project Structure

```
.
├── main.go              # Entry point — sets up routes and starts server
├── config/db.go         # Database connection
├── models/              # Data structs (User, Project)
├── controllers/         # Business logic for each route
├── routes/              # Route registration
├── init/001_schema.sql  # DB init script (runs automatically on first Docker start)
└── docker-compose.yml   # MariaDB service
```
