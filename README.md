# Public Library

Inspired by eminetto's [clean architecture](https://github.com/eminetto/clean-architecture-go-v2) backend API from 2020, this Public Library API project was developed from the ground up, utilizing hexagonal architecture. The purpose was to deepen my understanding of hexagonal architecture and provide a beginner-friendly project structure example. Although the project is not flawless, constructive feedback is highly appreciated to drive further improvements.

## What's hexagonal architecture

Hexagonal architecture, also known as Ports and Adapters architecture, is a software design pattern that promotes loose coupling and separation of concerns. It aims to isolate the core business logic of an application from the external dependencies and infrastructure details.

In hexagonal architecture, the core business logic resides at the center, surrounded by layers of adapters. These adapters act as interfaces between the core and the external systems such as databases, user interfaces, or third-party services. The core is oblivious to the existence of these adapters, which enables flexibility and ease of testing.

The architecture emphasizes the concept of ports and adapters. Ports define the interfaces through which the core interacts with the outside world, allowing it to remain independent of any specific technology or implementation details. Adapters, on the other hand, implement these interfaces and handle the communication between the core and the external systems.

By employing hexagonal architecture, developers can achieve better maintainability, testability, and flexibility in their applications. It helps in isolating business rules, decoupling from external dependencies, and allows for easier replacement or addition of components.

## Why you should use hexagonal architecture

- **Modularity** : Clear separation and isolation of components for improved maintainability.

- **Testability** : Core logic can be tested independently of external dependencies.

- **Flexibility** : Easy adaptation to changing requirements and integration with new technologies.

- **Decoupling** : Loose coupling between core logic and external systems for independent development.

- **Scalability** : Ability to scale different parts of the application independently.

- **Maintainability** : Localized changes and bug fixes without impacting the entire system.

- **Portability** : Core logic is independent of specific implementation details for easier migration.

## Why you shouldn't use hexagonal architecture

- **Simplicity** : For small or straightforward projects, a simpler architectural approach may be more suitable.

- **Overhead** : Hexagonal architecture introduces additional complexity and overhead, which may not be justified for resource-constrained projects.

- **Learning Curve** : Hexagonal architecture requires time and effort to learn and implement correctly, making it less suitable for time-sensitive or resource-limited projects.

- **Project Scope** : For small or short-lived projects, the benefits of hexagonal architecture may not outweigh the added effort and complexity.

- **Performance Considerations** : Hexagonal architecture's additional layers and abstractions can slightly impact performance, which may be a concern for performance-critical projects.

- **Project Constraints** : External factors, such as legacy systems or platform limitations, may not align well with hexagonal architecture.

## How to use this project

1. ### Configure the app.env File

    Start by configuring the app.env file to connect to your database. This file holds the necessary configuration settings for your application.

2. ### Verify Go Installation and Dependencies

    Before proceeding, ensure that you have Go installed on your system. Additionally, make sure you have installed all the necessary dependencies required by the project. Refer to the project documentation for information on installing dependencies.

### Run the project

```console
go run cmd/main.go
```

### Create book

```console
curl -X "POST" "http://localhost:8080/v1/books" \
-d $'{
    "title": "100 Go Mistakes and How to Avoid Them",
    "author": "Teiva Harsanyi",
    "amount": 5
}'
```

### Get book by ID

```console
curl -X "GET" "http://localhost:8080/v1/books/1"
```

### Search books by title

```console
curl -X "GET" "http://localhost:8080/v1/books" \
-H 'Accept: application/json' \
-d $'{ 
    "title": "100 Go Mistakes"
}'
```

### List books

```console
curl -X "GET" "http://localhost:8080/v1/books?"
```

### Update book

```console
curl -X "PUT" "http://localhost:8080/v1/books/1" \
-d $'{
    "title": "99 Go Mistakes and How to Avoid Them",
    "author": "Teiva Harsanyi",
    "amount": 5
}'
```

### Delete book

```console
curl -X "DELETE" "http://localhost:8080/v1/books/1"
```

### Create user

```console
curl -X "POST" "http://localhost:8080/v1/users" \
-d $'{
    "username": "Luigi",
    "password": "secret",
    "email": "luigi@email.com"
}'
```

### Get user

```console
curl -X "GET" "http://localhost:8080/v1/users/1"
```

### Update user

```console
curl -X "PUT" "http://localhost:8080/v1/users/1" \
-d $'{
    "username": "newUsername",
    "password": "superSecret",
    "email": "newemail@email.com"
}'
```

### Delete user

```console
curl -X "DELETE" "http://localhost:8080/v1/users/1"
```

### List all user loans

```console
curl -X "GET" "http://localhost:8080/v1/loans/1"
```

### Borrow book

```console
curl -X "POST" "http://localhost:8080/v1/loans/borrow" \
-d $'{
    "user_id": 1,
    "book_id": 1
}'
```

### Return book

```console
curl -X "POST" "http://localhost:8080/v1/loans/return" \
-d $'{
    "user_id": 1,
    "book_id": 1
}'
```

## Documentation

- [Database](https://dbdocs.io/luigiazevedo97/public_library_v2)

## Dependencies

- [Testify](https://github.com/stretchr/testify) - Tools for code testing

- [bcrypt](golang.org/x/crypto) - Password encryption function

- [pq](https://github.com/lib/pq) - Postgres driver

- [chi](https://github.com/go-chi/chi) - Lightweight, idiomatic and composable router for building Go HTTP services

- [zerolog](https://github.com/rs/zerolog) - A fast and simple logger

- [httplog](https://github.com/go-chi/httplog) - Small but powerful structured logging package for HTTP request logging in Go.

- [sqlmock](https://github.com/DATA-DOG/go-sqlmock) - Simulate any sql driver behavior in tests

- [Gomock](https://github.com/golang/mock) - Mocking framework

- [Viper](https://github.com/spf13/viper) - Viper is a complete configuration solution for Go applications.

## Tools used

- [Golang](https://go.dev/) - Programming language

- [task](https://github.com/go-task/task) - Task runner

- [Postgres](https://www.postgresql.org/) - Open source object-relational database

- [DBeaver](https://dbeaver.io/) - Database administration tool

- [Migrate](https://github.com/golang-migrate/migrate) - Database migrations

- [Golangci-lint](https://github.com/golangci/golangci-lint) - Automated code checking for programmatic and stylistic errors

- [dbdocs](https://dbdocs.io/?utm_source=dbdocs) - A free & simple tool to create web-based database documentation using DSL code.
