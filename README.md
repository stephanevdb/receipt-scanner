# Receipt Scanner API

This project is a backend service for a receipt scanning application, built with Go and extended with the PocketBase framework. It uses Google's Gemini AI to analyze uploaded receipt images and extract structured data.

## Features

- **Upload Receipts**: Upload receipt images to the server.
- **AI-Powered Analysis**: Uses Google's Gemini 1.5 Flash model to analyze receipts and extract information like store name, date, items, and total amount.
- **Data Persistence**: Stores receipt and item data in a SQLite database managed by PocketBase.
- **Verified Totals**: Automatically calculates the sum of item prices and verifies it against the receipt's total.
- **REST API**: Provides a comprehensive API for managing receipts and items.
- **Built-in Admin UI**: Leverages the PocketBase admin interface for easy data management.
- **Interactive API Documentation**: A user-friendly HTML page to explore and test the API endpoints.

## Tech Stack

- **Backend**: [Go](https://golang.org/)
- **Framework**: [PocketBase](https://pocketbase.io/)
- **AI**: [Google Gemini Pro](https://deepmind.google/technologies/gemini/)
- **Frontend (for documentation)**: HTML & [Tailwind CSS](https://tailwindcss.com/)

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) version 1.21 or higher.
- A Google Gemini API key. You can obtain one from [Google AI Studio](https://aistudio.google.com/app/apikey).

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/receipt-scanner.git
    cd receipt-scanner
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Create a `.env` file:**
    Create a `.env` file in the root of the project and add the following environment variables:

    ```env
    # Your Gemini API Key
    GEMINI_API_KEY="your_gemini_api_key"

    # Credentials for the initial PocketBase admin user
    POCKETBASE_INITIAL_ADMIN_EMAIL="admin@example.com"
    POCKETBASE_INITIAL_ADMIN_PASSWORD="your_secure_password"
    ```

### Running the Application

1.  **Start the server:**
    ```bash
    go run . serve
    ```

2.  The application will start, and you will see logs indicating that the collections have been initialized and the server is running. By default, it runs on `http://localhost:8080`.

## API Documentation

This project includes an interactive API documentation page. Once the server is running, you can access it at:

- **[http://localhost:8080/](http://localhost:8080/)**

This page provides a user-friendly interface to view all available endpoints, their descriptions, and example responses.

For a static, Markdown-formatted version of the API documentation, please see [API_DOCUMENTATION.md](./API_DOCUMENTATION.md).

### API Testing

A dedicated test page is available to try out the API endpoints in a web-based interface.

- **[http://localhost:8080/test](http://localhost:8080/test)**

From this page, you can perform health checks, upload receipts, trigger analysis, and list receipts and items.

### PocketBase Admin UI

You can access the PocketBase admin UI to view and manage your data directly:

- **[http://localhost:8080/_/](http://localhost:8080/_/)**

Log in with the admin credentials you set in your `.env` file.