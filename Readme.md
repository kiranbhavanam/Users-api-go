# User Management API

A **RESTful API** built with **Go** for managing user accounts with full CRUD operations and PostgreSQL database integration.

## 🚀 Features

- **Complete User Management**: Create, read, update, and delete user accounts
- **PostgreSQL Integration**: Persistent data storage with optimized queries
- **Input Validation**: Username, email, and data validation with business rules
- **Duplicate Prevention**: Prevents duplicate usernames and emails
- **Error Handling**: Custom error types with appropriate HTTP status codes
- **Clean Architecture**: Service layer, repository pattern, and dependency injection
- **Thread-Safe Operations**: Concurrent request handling

## 📋 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/users` | Retrieve all users |
| `GET` | `/users/{id}` | Get user by ID |
| `POST` | `/users` | Create new user |
| `PUT` | `/users/{id}` | Update existing user |
| `DELETE` | `/users/{id}` | Delete user |

## 🗃️ User Model

```json
{
    "id": 1,
    "username": "john_doe",
    "email": "john.doe@example.com",
    "password": "securepassword123",
    "name": "John Doe",
    "isactive": true
}
```

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **HTTP Router**: Gorilla Mux
- **Database**: PostgreSQL
- **Driver**: `lib/pq` (PostgreSQL driver)
- **Architecture**: Clean Architecture with Repository Pattern

## 🏃‍♂️ Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12+
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/user-management-api.git
   cd user-management-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL database**
   ```sql
   -- Create database
   CREATE DATABASE userdb;
   
   -- Connect to the database
   \c userdb;
   
   -- Create Users table
   CREATE TABLE Users (
       id SERIAL PRIMARY KEY,
       username VARCHAR(255) UNIQUE NOT NULL,
       email VARCHAR(255) UNIQUE NOT NULL,
       password VARCHAR(255) NOT NULL,
       name VARCHAR(255),
       isactive BOOLEAN DEFAULT true,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

4. **Configure database connection**
   
   Update the configuration in `internal/config/config.go`:
   ```go
   Host: "localhost"
   Port: 5433
   User: "postgres"
   Password: "password"
   DatabaseName: "userdb"
   SSLMode: "disable"
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

## 📁 Project Structure

```
user-management/
├── internal/
│   ├── config/         # Database configuration
│   ├── handlers/       # HTTP handlers and routing
│   ├── model/          # User data models
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic layer
│   └── errors/         # Custom error types
├── main.go            # Application entry point
├── go.mod             # Go module file
├── go.sum             # Go dependencies
└── README.md
```

## 🧪 API Usage Examples

### Get All Users
```bash
curl -X GET http://localhost:8080/users
```

### Get User by ID
```bash
curl -X GET http://localhost:8080/users/1
```

### Create New User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john.doe@example.com",
    "password": "securepassword123",
    "name": "John Doe"
  }'
```

### Update User
```bash
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johnsmith",
    "email": "john.smith@example.com",
    "password": "newpassword456",
    "name": "John Smith",
    "isactive": true
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:8080/users/1
```

## ✅ Validation Rules

### User Input Validation
- **Username**: Required, cannot be empty or same as name
- **Email**: Required, cannot be empty, must be unique
- **Name**: Required, cannot be empty
- **ID**: Must be positive (non-negative)

### Business Rules
- **Email Uniqueness**: Each email can only be associated with one user
- **Username Uniqueness**: Each username must be unique across all users
- **Update Validation**: Checks for duplicate username/email when updating existing users

## 🚦 HTTP Status Codes

| Status Code | Description |
|-------------|-------------|
| `200` | OK - Request successful |
| `201` | Created - User created successfully |
| `400` | Bad Request - Invalid input or validation error |
| `404` | Not Found - User not found |
| `409` | Conflict - Duplicate username or email |
| `500` | Internal Server Error - Database or server error |

## 🔍 Error Response Format

```json
{
  "error": "User with email 'john@example.com' already exists"
}
```

## 🔧 Configuration

### Database Configuration

The database configuration is located in `internal/config/config.go`:

```go
type DatabaseConfig struct {
    Host         string  // Database host
    Port         int     // Database port
    User         string  // Database user
    Password     string  // Database password
    DatabaseName string  // Database name
    SSLMode      string  // SSL mode (disable/require)
}
```

### Environment Setup

Make sure PostgreSQL is running on the configured port (default: 5433) and the database exists with the proper table structure.

## 🏗️ Architecture

The application follows **Clean Architecture** principles:

- **Handlers Layer**: HTTP request/response handling
- **Service Layer**: Business logic and validation
- **Repository Layer**: Data access and database operations
- **Model Layer**: Data structures and entities

### Key Components

- **PostgresRepository**: Handles all database operations
- **UserService**: Contains business logic and validation
- **UserHandler**: Manages HTTP requests and responses
- **Custom Errors**: Structured error handling with specific error types

## 🧪 Testing

### Manual Testing
Use the provided curl examples above or import the API collection into Postman for easier testing.

### Database Testing
Verify your database connection:
```bash
psql -h localhost -p 5433 -U postgres -d userdb -c "SELECT * FROM Users;"
```

## 🚧 Future Enhancements

- [ ] Password hashing with bcrypt
- [ ] JWT authentication and authorization
- [ ] Input sanitization for XSS prevention
- [ ] API rate limiting
- [ ] Request/response logging middleware
- [ ] Unit and integration tests
- [ ] Docker containerization
- [ ] API documentation with Swagger
- [ ] Environment-based configuration
- [ ] Database migrations

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

If you encounter any issues:
- Check the PostgreSQL connection and database setup
- Verify the table structure matches the expected schema
- Review the error logs for detailed error messages
- Create an issue on GitHub for additional support

***

**Note**: This is a development version. For production use, implement proper password hashing, authentication, and additional security measures.

[1](https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/attachments/86522524/931ead1d-da63-4375-bcb9-978abbf70453/main.go)
[2](https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/attachments/86522524/8c112055-45f9-4e42-8e07-feaecf9ca57c/user_repo.go)
[3](https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/attachments/86522524/acfb48bc-b38b-4461-931e-0dd3c3558970/user_service.go)
[4](https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/attachments/86522524/5f836604-df67-426d-81df-1db459eec803/config.go)
[5](https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/attachments/86522524/bc25f943-e262-4e87-a63c-c24621129a08/user_handlers.go)