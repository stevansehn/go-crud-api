# Go SQL Database API

A Go program that exposes a RESTful API for interacting with a MySQL database, allowing for CRUD operations on a `users` table through function handlers. This application enables you to insert new users, retrieve user information, update existing users, and delete users from the database, all through HTTP requests.

## Features

- RESTful API for user management.
- Connects to a MySQL database.
- Supports CRUD operations: Create, Read, Update, and Delete users.
- Utilizes `github.com/gorilla/mux` for routing and handling HTTP requests.
- Returns structured JSON responses for easy integration with front-end applications.
- Ensures secure password handling (consider hashing in future updates).

## Prerequisites

- Go (version 1.16 or later)
- MySQL server
- Go MySQL driver: `github.com/go-sql-driver/mysql`
- Gorilla Mux router: `github.com/gorilla/mux`
- Environment variable loader: `github.com/joho/godotenv`

## Installation

1. **Clone the repository:**
   ```bash
   git clone git@github.com:yourusername/go-crud-api.git
   cd go-sql-database-api
   ```

2. **Install dependencies:**
   Make sure you have Go installed, then run:
   ```bash
   go mod tidy
   ```

3. **Set up your MySQL database:**
   - Create a new database (e.g., `go_sql_database_api`).
   - Update the `.env` file with your database credentials:
     ```
     DB_USER=your_db_user
     DB_PASSWORD=your_db_password
     DB_HOST=localhost
     DB_PORT=3306
     DB_NAME=go_sql_database_api
     DB_PARSE_TIME=true
     ```

4. **Create the `users` table:**
   The application will automatically create the `users` table if it does not exist when you run the program.

## Usage

1. **Run the application:**
   ```bash
   go run main.go
   ```

2. **API Endpoints:**
   - **Create User**: `POST /users`
   - **Get User**: `GET /users/{id}`
   - **Update User**: `PUT /users/{id}`
   - **Delete User**: `DELETE /users/{id}`
   - **Get All Users**: `GET /users`

3. **Expected Output:**
   Upon making requests to the API, you should receive structured JSON responses indicating the success of operations and user data.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
