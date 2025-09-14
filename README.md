# CSV Sales API

A well-structured Go backend API for processing CSV sales data. This API accepts CSV file uploads, aggregates sales data by department, and returns downloadable result files.

## Features

- **CSV File Upload**: Accept CSV files via HTTP POST endpoint
- **Sales Aggregation**: Automatically aggregate total sales per department
- **Streaming Processing**: Memory-efficient processing of large CSV files using streaming
- **File Management**: Save uploaded files and generated results with UUID-based naming
- **Download Links**: Return downloadable links to processed result files
- **Error Handling**: Comprehensive error handling and logging
- **Modular Design**: Clean separation of concerns with proper Go project structure
- **Unit Tests**: Comprehensive test coverage for core functionality
- **CORS Support**: Built-in CORS middleware for web applications

## File Processing Complexity

The CSV file processing logic is highly efficient:

- **Time Complexity:** O(n), where n is the number of data rows in the CSV file. Each row is read and processed once.
- **Space Complexity:** O(m), where m is the number of unique departments (for aggregation).

This ensures the backend can handle large CSV files with minimal memory usage and fast processing.

## Project Structure

```
csv-sales-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/
│   │   └── upload_handler.go    # HTTP request handlers
│   ├── models/
│   │   └── sales.go            # Data models
│   └── services/
│       ├── csv_service.go      # CSV processing logic
│       ├── csv_service_test.go # CSV service tests
│       ├── file_service.go     # File operations
│       └── file_service_test.go # File service tests
├── pkg/
│   └── utils/
│       └── env.go              # Utility functions
├── public/
│   └── uploads/                # File storage directory
├── examples/                   # Sample CSV files
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── .gitignore                  # Git ignore rules
└── README.md                   # This file
```

## Prerequisites

- Go 1.21 or higher
- Git

## Installation & Setup

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd csv-sales-api
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Create uploads directory** (if not exists):
   ```bash
   mkdir -p public/uploads
   ```

4. **Run the application**:
   ```bash
   go run cmd/server/main.go
   ```

   The server will start on `http://localhost:8080` by default.

## Usage

### Upload and Process CSV

**Endpoint**: `POST /api/v1/upload`

**Request**: Multipart form data with CSV file

**Example using curl**:
```bash
curl -X POST \
  -F "file=@examples/sample_sales.csv" \
  http://localhost:8080/api/v1/upload
```

**Response**:
```json
{
  "success": true,
  "message": "CSV file processed successfully",
  "download_url": "/public/uploads/result_12345678-1234-1234-1234-123456789abc.csv",
  "total_departments": 4,
  "processed_at": "2024-01-15T10:30:00Z"
}
```

### Health Check

**Endpoint**: `GET /api/v1/health`

**Response**:
```json
{
  "status": "ok"
}
```

### Download Result File

Access the result file directly via the download URL:
```
http://localhost:8080/public/uploads/result_12345678-1234-1234-1234-123456789abc.csv
```

The result CSV file will contain two columns:
- **Department Name**: The name of each department
- **Total Number of Sales**: The aggregated sales total for that department

## CSV Format Requirements

The API expects CSV files with the following characteristics:

- **Required Columns**: Must contain columns for department and sales data
- **Flexible Headers**: Supports various column names:
  - Department: `department`, `dept`
  - Sales: `sales`, `total_sales`, `total sales`, `amount`, `revenue`
- **File Size**: No limit (optimized for large files)
- **File Type**: Only `.csv` files are accepted

### Example CSV Format

```csv
department,sales
Electronics,1500
Clothing,800
Books,300
Electronics,2000
Home & Garden,1200
```

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Test Coverage

The project includes comprehensive unit tests for:
- CSV processing logic
- File validation and operations
- Error handling scenarios
- Edge cases and malformed data

## API Error Handling

The API returns structured error responses:

```json
{
  "success": false,
  "error": "Error description",
  "code": 400
}
```

### Common Error Codes

- `400`: Bad Request (invalid file, missing file, validation errors)
- `500`: Internal Server Error (processing failures, file system errors)
