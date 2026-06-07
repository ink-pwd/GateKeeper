# Gatekeeper

Gatekeeper is a lightweight authentication microservice written in Go.

The service provides user registration, authentication, JWT token generation, and token validation through a REST API. User data is stored in PostgreSQL, while passwords are securely hashed using bcrypt.

## Features

* User registration
* User authentication
* JWT token generation
* JWT token validation
* Password hashing with bcrypt
* PostgreSQL integration
* Input validation
* Environment-based configuration
* Structured logging

## API

### POST /register

Registers a new user.

### POST /login

Authenticates a user and returns a JWT token.

### GET /validate

Validates a JWT token and returns the associated user email.

## Tech Stack

* Go
* PostgreSQL
* JWT
* bcrypt
* REST API

## Project Goal

The purpose of this project is to demonstrate the implementation of a simple authentication service using Go, including secure password storage, token-based authentication, database integration, and clean project structure.
