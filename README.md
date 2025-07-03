# Todo App - Golang

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-v2.52.8-00A3E0?style=flat&logo=go)](https://gofiber.io/)
[![GORM](https://img.shields.io/badge/GORM-Latest-00ADD8?style=flat&logo=go)](https://gorm.io/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A robust RESTful API for managing todo tasks built with Go, Fiber, and GORM. This application provides comprehensive task management capabilities with pagination, filtering, and search functionality.

## 🌐 Live Demo

Explore the live demo: [https://todo-app-golang.jefripunza.com](https://todo-app-golang.jefripunza.com)

## ✨ Features

- **Create Tasks**: Add new tasks with title, description, status, and due date
- **List Tasks**: Get all tasks with pagination, filtering, and search capabilities
- **Get Task Details**: Retrieve specific task information by ID
- **Update Tasks**: Modify existing task details
- **Delete Tasks**: Remove tasks from the system
- **API Documentation**: Comprehensive Swagger documentation

## 🛠️ Tech Stack

- **Go**: Programming language
- **Fiber**: Web framework
- **GORM**: ORM library for database operations
- **PostgreSQL**: Database
- **Swagger**: API documentation

## 📋 API Documentation

### Base URL

```
/
```

### Endpoints

#### Tasks

| Method | Endpoint     | Description                                  |
|--------|-------------|----------------------------------------------|
| GET    | /tasks      | List all tasks with pagination and filtering |
| POST   | /tasks      | Create a new task                           |
| GET    | /tasks/{id} | Get a specific task by ID                    |
| PUT    | /tasks/{id} | Update a specific task                       |
| DELETE | /tasks/{id} | Delete a specific task                       |

### Query Parameters for GET /tasks

| Parameter | Type   | Required | Description                                       |
|-----------|--------|----------|---------------------------------------------------|
| status    | string | No       | Filter tasks by status (pending/completed)        |
| page      | string | No       | Page number for pagination (default: 1)           |
| limit     | string | No       | Number of tasks per page (default: 10)            |
| search    | string | No       | Search term to filter by title or description     |

### Request Bodies

#### Create Task (POST /tasks)

```json
{
  "title": "Task title",
  "description": "Task description",
  "status": "pending",
  "due_date": "2025-07-02 15:04:05"
}
```

#### Update Task (PUT /tasks/{id})

```json
{
  "title": "Updated task title",
  "description": "Updated task description",
  "status": "completed",
  "due_date": "2025-07-05 22:19:00"
}
```

### Response Objects

#### Todo Object

```json
{
  "id": 1,
  "title": "Task title",
  "description": "Task description",
  "status": "pending",
  "due_date": "2025-07-02T15:04:05Z"
}
```

#### Paginated Response

```json
{
  "tasks": [
    {
      "id": 1,
      "title": "Task title",
      "description": "Task description",
      "status": "pending",
      "due_date": "2025-07-02T15:04:05Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_tasks": 42
  }
}
```

#### Error Response

```json
{
  "code": "ERROR_CODE",
  "message": "Error message"
}
```

## 🚀 Getting Started

### Prerequisites

- Go 1.20 or higher
- PostgreSQL

### Installation

1. Clone the repository
   ```bash
   git clone https://github.com/jefripunza/todo-app-golang.git
   cd todo-app-golang
   ```

2. Install dependencies
   ```bash
   go mod download
   ```

3. Configure environment variables
   ```bash
   cp .env.example .env
   # Edit .env file with your database credentials
   ```

4. Run the application
   ```bash
   go run main.go
   ```

5. Access the API at `http://localhost:3000`

### Docker

```bash
# Build the Docker image
docker build -t todo-app-golang .

# Run the container
docker run -p 3000:3000 todo-app-golang
```

## 📝 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Jefri Herdi Triyanto**

- Website: [jefripunza.com](https://jefripunza.com)
- Email: [hi@jefripunza.com](mailto:hi@jefripunza.com)
