package models

// SalesRecord represents a single sales record from CSV
type SalesRecord struct {
	Department string `csv:"department"`
	Sales      int    `csv:"sales"`
}

// DepartmentSummary represents aggregated sales data for a department
type DepartmentSummary struct {
	Department string `json:"department" csv:"Department Name"`
	TotalSales int    `json:"total_sales" csv:"Total Number of Sales"`
}

// UploadResponse represents the response after successful CSV upload and processing
type UploadResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	DownloadURL      string `json:"download_url"`
	TotalDepartments int    `json:"total_departments"`
	TotalSales       int    `json:"total_sales"`
	ProcessedAt      string `json:"processed_at"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}
