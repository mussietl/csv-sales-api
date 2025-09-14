package services

import (
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileServiceValidateFile(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_uploads")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	fileService := NewFileService(tempDir, logger)

	tests := []struct {
		name        string
		fileHeader  *multipart.FileHeader
		expectError bool
	}{
		{
			name: "valid CSV file",
			fileHeader: &multipart.FileHeader{
				Filename: "test.csv",
				Size:     1000,
				Header:   map[string][]string{"Content-Type": {"text/csv"}},
			},
			expectError: false,
		},
		{
			name: "large file (no size limit)",
			fileHeader: &multipart.FileHeader{
				Filename: "test.csv",
				Size:     11 << 20, // 11MB
				Header:   map[string][]string{"Content-Type": {"text/csv"}},
			},
			expectError: false,
		},
		{
			name: "invalid file extension",
			fileHeader: &multipart.FileHeader{
				Filename: "test.txt",
				Size:     1000,
				Header:   map[string][]string{"Content-Type": {"text/plain"}},
			},
			expectError: true,
		},
		{
			name: "valid CSV with different MIME type",
			fileHeader: &multipart.FileHeader{
				Filename: "test.csv",
				Size:     1000,
				Header:   map[string][]string{"Content-Type": {"application/csv"}},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fileService.ValidateFile(tt.fileHeader)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFileServiceSaveResultFile(t *testing.T) {
	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_uploads")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	fileService := NewFileService(tempDir, logger)

	// Test data
	summaries := []DepartmentSummary{
		{Department: "Electronics", TotalSales: 2500},
		{Department: "Clothing", TotalSales: 700},
		{Department: "Books", TotalSales: 300},
	}

	// Save result file
	filePath, err := fileService.SaveResultFile(summaries)
	require.NoError(t, err)
	assert.NotEmpty(t, filePath)

	// Verify file exists
	assert.FileExists(t, filePath)

	// Read and verify file contents
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	expectedContent := "Department Name,Total Number of Sales\nElectronics,2500\nClothing,700\nBooks,300\n"
	assert.Equal(t, expectedContent, string(content))
}

func TestFileServiceGetDownloadURL(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_uploads")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := logrus.New()
	fileService := NewFileService(tempDir, logger)

	// Test with full path
	fullPath := filepath.Join(tempDir, "test_file.csv")
	url := fileService.GetDownloadURL(fullPath)
	expectedURL := "/public/uploads/test_file.csv"
	assert.Equal(t, expectedURL, url)

	// Test with just filename
	url = fileService.GetDownloadURL("test_file.csv")
	assert.Equal(t, expectedURL, url)
}

func TestFileServiceSaveUploadedFile(t *testing.T) {
	// This test is simplified since we can't easily mock multipart.FileHeader.File()
	// In a real scenario, we would use dependency injection or interfaces
	// For now, we'll test the other functionality

	// Create temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test_uploads")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := logrus.New()
	logger.SetLevel(logrus.ErrorLevel)
	fileService := NewFileService(tempDir, logger)

	// Test that the service is properly initialized
	assert.NotNil(t, fileService)
	assert.Equal(t, tempDir, fileService.uploadsDir)
	assert.NotNil(t, fileService.logger)
}
