# Backend Engineer Assessment Test

## Tech Stack
- **Programming Language**: Go
- **Web Framework**: Echo
- **Database**: PostgreSQL
- **Containerization**: Docker & Docker Compose

## Project Structure
```
- config/
  - config.go        # Configuration setup
- database/
  - database.go      # Database connection and setup
- handler/
  - nasabah.go       # Business logic for handling customer operations
- model/
  - nasabah.go       # Data models
- middleware/
  - logger.go        # Logging middleware
- router/
  - router.go        # API route definitions
- .env               # Environment variables
- docker-compose.yaml  # Docker Compose configuration
- Dockerfile         # Docker build file
- main.go            # Entry point of the application
```

## API Endpoints
### 1. **Register Nasabah**
**Endpoint:** `POST /daftar`

**Request Body:**
```json
{
  "nama": "John Doe",
  "nik": "1234567890123456",
  "no_hp": "081234567890"
}
```

**Response:**
```json
{
  "no_rekening": "1234567890"
}
```

**Error Responses:**
- `400`: Jika NIK atau nomor HP sudah terdaftar

---
### 2. **Deposit (Tabung)**
**Endpoint:** `POST /tabung`

**Request Body:**
```json
{
  "no_rekening": "1234567890",
  "nominal": 500000
}
```

**Response:**
```json
{
  "saldo": 1500000
}
```

**Error Responses:**
- `400`: Jika nomor rekening tidak ditemukan

---
### 3. **Withdraw (Tarik)**
**Endpoint:** `POST /tarik`

**Request Body:**
```json
{
  "no_rekening": "1234567890",
  "nominal": 200000
}
```

**Response:**
```json
{
  "saldo": 1300000
}
```

**Error Responses:**
- `400`: Jika nomor rekening tidak ditemukan atau saldo tidak cukup

---
### 4. **Check Balance (Saldo)**
**Endpoint:** `GET /saldo/:no_rekening`

**Response:**
```json
{
  "saldo": 1300000
}
```

**Error Responses:**
- `400`: Jika nomor rekening tidak ditemukan

---

## Deployment & Setup

### 1. **Environment Variables (.env file)**
```
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=banking
DB_PORT=5432
```

### 2. **Run with Docker**
```sh
docker-compose up --build
```

### 3. **Accessing the PostgreSQL Database**
If using **DBeaver**, connect using the following details:
```
Host: localhost
Port: 5432
Database: banking
Username: postgres
Password: postgres
```

## Logging
The application uses structured logging with appropriate log levels:
- **INFO**: Normal operations (e.g., user registration, transactions, etc.)
- **WARN**: Invalid inputs or failed validations
- **ERROR**: System errors (e.g., database failures)

## Assessment Criteria
- **Logging**: Structured and informative logs
- **Software Architecture**: Clean module separation and naming conventions
- **Database Schema**: Properly structured relationships
- **Database Operations**: Efficient queries
- **REST API Design**: Well-defined endpoints, parameters, and middleware
- **Configuration Management**: Secure and flexible configurations
- **Deployment**: Dockerized application setup

## Repository & Video
- **GitHub Repository**: [gobanking](https://github.com/restuedos/gobanking)
- **Demo Video**: [Watch Here](https://drive.google.com/file/d/130Tva431DdR7692aIlgLWBA2y1bOpleM/view?usp=sharing)

---
Developed by **Restu Edo Setiaji**

