# Receipt Scanner API Documentation

This document provides a detailed overview of the API endpoints for the Receipt Scanner application.

---

## Health Check

- **Endpoint:** `GET /api/health`
- **Description:** Returns the health status of the API server.

**Success Response (`200 OK`):**
```json
{
  "status": "ok"
}
```

---

## Upload Receipt

- **Endpoint:** `POST /api/receipts/upload`
- **Description:** Uploads a receipt image. The request must be a `multipart/form-data` POST request with the image attached to a form field named `receipt`.

**Success Response (`200 OK`):**
```json
{
  "message": "Receipt uploaded successfully.",
  "filename": "receipt-2023-01.jpg"
}
```

---

## Analyze Receipt

- **Endpoint:** `POST /api/receipts/analyze`
- **Description:** Analyzes a previously uploaded receipt image using AI.

**Request Body:**
```json
{
  "filename": "receipt-2023-01.jpg"
}
```

**Success Response (`200 OK`):**
```json
{
  "title": "Grocery Store",
  "date": "2023-10-27",
  "items": [
    { "name": "Item 1", "price": 10.99 },
    { "name": "Item 2", "price": 5.49 }
  ],
  "total": 16.48,
  "verified_total": true
}
```

---

## List All Receipts

- **Endpoint:** `GET /api/receipts`
- **Description:** Lists all uploaded and analyzed receipts.

**Success Response (`200 OK`):**
```json
[
  {
    "id": "RECORD_ID",
    "created": "2024-01-01T12:00:00.000Z",
    "filename": "receipt1.jpg",
    "title": "Grocery Store",
    "total": 25.50,
    "verified_total": true
  }
]
```

---

## List Items in a Receipt

- **Endpoint:** `GET /api/receipts/:id/items`
- **Description:** Lists all items for a specific receipt.

**Success Response (`200 OK`):**
```json
[
  {
    "name": "Item 1",
    "price": 10.99
  }
]
``` 