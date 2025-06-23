# Receipt Scanner API

This project is a backend service for a receipt scanning application, built with Go and extended with the PocketBase framework. It uses Google's Gemini AI to analyze uploaded receipt images and extract structured data.

## Features

- **Upload Receipts**: Upload receipt images to the server.
- **AI-Powered Analysis**: Uses Google's Gemini 1.5 Flash model to analyze receipts and extract information like store name, date, items, and total amount.
- **Data Persistence**: Stores receipt and item data in a SQLite database managed by PocketBase.
- **Verified Totals**: Automatically calculates the sum of item prices and verifies it against the receipt's total.
- **User Management**: Basic user registration.
- **REST API**: Provides a comprehensive API for managing receipts, items, and users.
- **Built-in Admin UI**: Leverages the PocketBase admin interface for easy data management.
- **Interactive API Documentation**: A user-friendly HTML page to explore and test the API endpoints.

## Tech Stack

- **Backend**: [Go](https://golang.org/)
- **Framework**: [PocketBase](https://pocketbase.io/)
- **AI**: [Google Gemini 1.5 Flash](https://deepmind.google/technologies/gemini/)
- **Frontend (for documentation)**: HTML & [Tailwind CSS](https://tailwindcss.com/)

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) version 1.21 or higher.
- A Google Gemini API key. You can obtain one from [Google AI Studio](https://aistudio.google.com/app/apikey).

### Installation & Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your-username/receipt-scanner-go.git
    cd receipt-scanner-go
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
    go run .
    ```

2.  The application will start, and you will see logs indicating that the collections have been initialized and the server is running. By default, it runs on `http://localhost:8090`.

## API Endpoints

An interactive API documentation page is available at [http://localhost:8090/](http://localhost:8090/).

For a quick reference, here are the available endpoints:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/api/health` | Checks the health of the application. |
| `GET` | `/api/receipts` | Retrieves a list of all receipts. |
| `POST` | `/api/receipts/upload` | Uploads a receipt image. |
| `POST` | `/api/receipts/analyze` | Analyzes an uploaded receipt and saves the data. |
| `GET` | `/api/receipts/:id/items` | Retrieves all items for a specific receipt. |
| `DELETE`| `/api/receipts/:id` | Deletes a receipt and all its associated items. |
| `GET` | `/api/items/:id` | Retrieves a specific item by its ID. |
| `PATCH` | `/api/items/:id/paid` | Updates the paid status of a specific item. |
| `POST` | `/api/users/create` | Creates a new user. |

### API Testing

A dedicated test page is available to try out the API endpoints in a web-based interface.

- **[http://localhost:8090/test](http://localhost:8090/test)**

From this page, you can perform health checks, upload receipts, trigger analysis, and list receipts and items.

### PocketBase Admin UI

You can access the PocketBase admin UI to view and manage your data directly:

- **[http://localhost:8090/_/](http://localhost:8090/_/)**

Log in with the admin credentials you set in your `.env` file.

## Database Schema

The application uses three collections in its PocketBase database:

### `receipts` Collection

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `filename` | Text | Yes | The name of the uploaded receipt file. |
| `title` | Text | No | The name of the store. |
| `date` | Date | No | The date of the receipt. |
| `total` | Number | No | The total amount from the receipt. |
| `verified_total`| Bool | No | True if the sum of items matches the total. |

### `items` Collection

| Field | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | Text | Yes | The name of the item. |
| `price` | Number | Yes | The price of a single item. |
| `quantity` | Number | No | The quantity of the item (defaults to 1). |
| `amount` | Number | No | The total amount for the item (`price` * `quantity`). |
| `paid` | Number | No | The quantity of the item that has been paid for. |
| `receipt`| Relation | Yes | A link to the parent receipt. |

### `users` Collection

This is the standard PocketBase `users` collection, with an added custom field.

| Field | Type | Description |
| --- | --- | --- |
| `name` | Text | The user's name. |