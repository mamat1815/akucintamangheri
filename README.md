# Campus Lost & Found API

Backend service for the Campus Lost & Found system, designed with privacy-first principles and a smart matching engine.

## 🚀 Tech Stack

- **Language:** Go (Golang) 1.22+
- **Framework:** Gin Gonic
- **Database:** PostgreSQL (via Supabase)
- **ORM:** GORM
- **Authentication:** JWT (JSON Web Tokens)
- **Documentation:** Swagger (Swaggo)

## 🛠️ Prerequisites

- Go 1.22 or higher
- PostgreSQL database (or Supabase project)

## ⚙️ Setup & Installation

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
    Create a `.env` file in the root directory based on the example below:

    ```env
    DATABASE_URL="postgresql://user:password@host:port/dbname?sslmode=disable"
    JWT_SECRET="your_super_secret_key"
    PORT=8080
    ```

4.  **Run the Server**
    ```bash
    go run cmd/server/main.go
    ```
    The server will start on port `8080` (or the port defined in `.env`). Database migrations will run automatically on startup.

## 📚 API Documentation

Once the server is running, you can access the interactive Swagger API documentation at:

👉 **http://localhost:8080/swagger/index.html**

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
