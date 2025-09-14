package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// FileService handles file operations
type FileService struct {
	uploadsDir string
	logger     *logrus.Logger
}

// NewFileService creates a new FileService instance
func NewFileService(uploadsDir string, logger *logrus.Logger) *FileService {
	return &FileService{
		uploadsDir: uploadsDir,
		logger:     logger,
	}
}

// SaveUploadedFile saves an uploaded file to the uploads directory
func (fs *FileService) SaveUploadedFile(file *multipart.FileHeader) (string, error) {
	// Generate unique filename
	fileExt := filepath.Ext(file.Filename)
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("upload_%s%s", uniqueID, fileExt)
	filePath := filepath.Join(fs.uploadsDir, filename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		fs.logger.Errorf("Failed to open uploaded file: %v", err)
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		fs.logger.Errorf("Failed to create destination file: %v", err)
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, src); err != nil {
		fs.logger.Errorf("Failed to copy file contents: %v", err)
		return "", fmt.Errorf("failed to copy file contents: %w", err)
	}

	fs.logger.Infof("File saved successfully: %s", filePath)
	return filePath, nil
}

// SaveResultFile saves the aggregated results to a CSV file
func (fs *FileService) SaveResultFile(departmentSummaries []DepartmentSummary) (string, error) {
	// Generate unique filename for result
	uniqueID := uuid.New().String()
	filename := fmt.Sprintf("result_%s.csv", uniqueID)
	filePath := filepath.Join(fs.uploadsDir, filename)

	// Create result file
	file, err := os.Create(filePath)
	if err != nil {
		fs.logger.Errorf("Failed to create result file: %v", err)
		return "", fmt.Errorf("failed to create result file: %w", err)
	}
	defer file.Close()

	// Write CSV header
	if _, err := file.WriteString("Department Name,Total Number of Sales\n"); err != nil {
		fs.logger.Errorf("Failed to write CSV header: %v", err)
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, summary := range departmentSummaries {
		line := fmt.Sprintf("%s,%d\n", summary.Department, summary.TotalSales)
		if _, err := file.WriteString(line); err != nil {
			fs.logger.Errorf("Failed to write CSV data: %v", err)
			return "", fmt.Errorf("failed to write CSV data: %w", err)
		}
	}

	fs.logger.Infof("Result file saved successfully: %s", filePath)
	return filePath, nil
}

// GetDownloadURL generates a download URL for a file
func (fs *FileService) GetDownloadURL(filePath string) string {
	// Extract just the filename from the full path
	filename := filepath.Base(filePath)
	return fmt.Sprintf("/public/uploads/%s", filename)
}

// ValidateFile validates the uploaded file
func (fs *FileService) ValidateFile(file *multipart.FileHeader) error {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".csv" {
		return fmt.Errorf("only CSV files are allowed, got: %s", ext)
	}

	// Check MIME type
	if !strings.Contains(file.Header.Get("Content-Type"), "text/csv") &&
		!strings.Contains(file.Header.Get("Content-Type"), "application/csv") &&
		!strings.Contains(file.Header.Get("Content-Type"), "text/plain") {
		fs.logger.Warnf("Unexpected MIME type: %s", file.Header.Get("Content-Type"))
		// Don't fail here as some systems may not set the correct MIME type
	}

	fs.logger.Infof("File validation passed for file: %s (size: %d bytes)", file.Filename, file.Size)
	return nil
}

// DepartmentSummary represents aggregated sales data for a department
type DepartmentSummary struct {
	Department string `json:"department" csv:"Department Name"`
	TotalSales int    `json:"total_sales" csv:"Total Number of Sales"`
}
