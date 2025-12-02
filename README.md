# Campus Lost & Found API

Backend service for the Campus Lost & Found system, designed with privacy-first principles and a smart matching engine.

## 🚀 Tech Stack

- **Language:** Go (Golang) 1.24
- **Framework:** Gin Gonic
- **Database:** PostgreSQL (via Supabase)
- **ORM:** GORM
- **Authentication:** JWT (JSON Web Tokens)
- **Documentation:** Swagger (Swaggo)
- **Deployment:** Docker & Jenkins (CI/CD)

## 🛠️ Prerequisites

- Go 1.24 or higher
- Docker & Docker Compose (optional for local dev)
- PostgreSQL database (or Supabase project)

## ⚙️ Setup & Installation

### Local Development

1.  **Clone the repository**
    ```bash
    git clone <repository-url>
    cd akucintamangheri
    ```

2.  **Install Dependencies**
    ```bash
    go mod tidy
    ```

3.  **Environment Variables**
    Create a `.env` file in the root directory. You can copy the example below:

    ```env
    # Database
    DATABASE_URL="postgresql://user:password@host:port/dbname?sslmode=disable"
    
    # Server
    PORT=3000
    GIN_MODE=debug # or release
    
    # Auth
    JWT_SECRET="your_super_secret_key"
    JWT_EXPIRY=24h
    
    # CORS
    ALLOWED_ORIGINS=http://localhost:4200,http://localhost:8080
    
    # File Upload
    MAX_UPLOAD_SIZE=10485760 # 10MB
    UPLOAD_PATH=./uploads
    ```

4.  **Run the Server**
    ```bash
    go run cmd/server/main.go
    ```
    The server will start on port `3000` (default).

### Docker Deployment (VPS)

The project is configured for automated deployment via Jenkins using Docker.

1.  **Dockerfile**: Builds a lightweight Alpine-based image.
2.  **Jenkinsfile**: Automates the build and deploy process.
    -   Builds Docker image `campus-lost-found`.
    -   Runs container `campus-lost-found-container`.
    -   Mounts `.env` from `/var/www/campus-lost-found/.env`.
    -   Mounts storage from `/var/www/campus-lost-found/storage`.

**VPS Requirements:**
-   Docker installed & user added to `docker` group.
-   Directory `/var/www/campus-lost-found` created.
-   `.env` file present in that directory.

## 📚 API Documentation

### Live Swagger Docs (VPS)
👉 **http://api.afsar.my.id/swagger/index.html**

### Local Swagger Docs
👉 **http://localhost:3000/swagger/index.html**

To update Swagger docs after code changes:
```bash
swag init -g cmd/server/main.go --output docs
```

## ✨ Key Features

-   **Authentication**: User registration and login with Role-Based Access Control (USER, ADMIN, SECURITY).
-   **Asset Management (Owner-First)**: Register assets with private details. Generates unique QR codes for each asset.
-   **Privacy Mode**: Scanning a QR code reveals only public info (Category, Nearest Security Point) unless the owner views it.
-   **Lost & Found Workflow**:
    -   **Report Lost**: Owners can mark assets as lost.
    -   **Report Found (QR)**: Finders scan QR to report location.
    -   **Report Found (No QR)**: Finders report items they found.
-   **Smart Matching Engine**: Automatically matches "Found Items" (without QR) to "Lost Assets" based on category and time.
-   **Claims System**: Owners can claim found items by answering verification questions.
-   **Notifications**: In-app notifications for matches and claim updates.
-   **File Uploads**: Secure image uploads for assets and found items.

## 📂 Project Structure

-   `cmd/server`: Application entry point.
-   `config`: Configuration and DB connection.
-   `internal/controllers`: HTTP Request handlers.
-   `internal/services`: Business logic.
-   `internal/repository`: Database interactions.
-   `internal/models`: Database schemas.
-   `internal/dto`: Data Transfer Objects for API I/O.
-   `internal/middleware`: Auth and CORS middleware.
-   `internal/matching`: Smart matching logic.
-   `docs`: Swagger documentation files.
