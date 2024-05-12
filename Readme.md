# Tigerhall Kittens API

A GraphQL API built with Go to track tiger sightings in the wild.

## Table of Contents

-   [Project Overview](#project-overview)
-   [Features](#features)
-   [Technologies Used](#technologies-used)
-   [Getting Started](#getting-started)
    -   [Prerequisites](#prerequisites)
    -   [Installation](#installation)
    -   [Database Setup](#database-setup)
    -   [Running the Server](#running-the-server)
    -   [Running Tests](#running-tests)
-   [API Documentation](#api-documentation)
    -   [Authentication](#authentication)
    -   [Mutations](#mutations)
    -   [Queries](#queries)
-   [Error Handling](#error-handling)
-   [Additional Notes](#additional-notes)

## Project Overview

This project is a backend API designed to support a fictional mobile app that allows users to report tiger sightings. The API exposes functionality for user registration and login, creation of tiger profiles, and recording new sightings with image uploads. It also includes a notification system to alert users who have previously reported sightings of the same tiger.

## Features

*   **User Management:** Registration and login with JWT authentication.
*   **Tiger Management:** Creation and listing of tiger profiles (with pagination).
*   **Sighting Management:** Creation of tiger sightings with location, timestamp, and image upload.
*   **Distance Restriction:** Enforces a 5km distance rule for new sightings of the same tiger.
*   **Notifications:**  Alerts users who have previously sighted the same tiger when a new sighting is reported (implementation using Go channels or an external message queue).
*   **Error Handling:** Provides informative error messages and appropriate HTTP status codes.
*   **Image Resizing:** Resizes uploaded images using the `imaging` library.

## Technologies Used

*   **Language:** Go (Golang)
*   **Web Framework:** Gin
*   **GraphQL:** gqlgen
*   **Database:** PostgreSQL
*   **ORM:** GORM
*   **Authentication:** JWT (JSON Web Tokens)
*   **Image Resizing:** github.com/disintegration/imaging
*   **Logging:** logrus

## Getting Started

### Prerequisites

*   Go (v1.21.4 or later) 
*   PostgreSQL
*   (Optional) Mailtrap or a similar service for testing email notifications

### Installation

1.  Clone the repository:
```bash
git clone https://github.com/nurcholisnanda/tigerhall-kittens.git
cd tigerhall-kittens
```
2. Update dependency
```bash
go mod tidy
```
3. Generate graphql (if necessary)
```bash
go generate ./...
```
4. Build and run with docker-compose
```bash
docker-compose build
docker-compose up
```
5. Running Tests
```bash
go test ./...
```

## API Documentation

Detailed API documentation (queries, mutations, input types) can be found in the GraphQL Playground (or similar URL:"localhost:8080") after starting the server.

### Authentication

*   Use the `login` mutation to obtain a JWT token.
*   Include the token in the `Authorization` header for mutation requests:

## Error Handling

The API uses a structured error format with `message` and `extensions` fields. Error codes are included in the `extensions` to provide additional information to the client.

## Additional Notes

*   The API currently uses Go channels for notifications. For production, consider using a more robust message queue system (e.g., RabbitMQ, Kafka).
*   Remember to replace placeholders (like database credentials) with your actual configuration values.

